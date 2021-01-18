package crds

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	pinot_crds "github.com/spaghettifunk/pinot-operator/pkg/manifests/pinot-crds/generated"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
)

const (
	componentName     = "crds"
	createdByLabel    = "apache.io/created-by"
	createdBy         = "pinot-operator"
	eventRecorderName = "pinot-crd-controller"
)

type CRDReconciler struct {
	crds     []runtime.Object
	config   *rest.Config
	revision string
	recorder record.EventRecorder
	client   client.Client
}

func New(mgr manager.Manager, revision string, crds ...runtime.Object) (*CRDReconciler, error) {
	r := &CRDReconciler{
		crds:     crds,
		config:   mgr.GetConfig(),
		revision: revision,
		recorder: mgr.GetEventRecorderFor(eventRecorderName),
		client:   mgr.GetClient(),
	}

	return r, nil
}

func (r *CRDReconciler) LoadCRDs() error {
	dir, err := pinot_crds.CRDs.Open("/")
	if err != nil {
		return err
	}

	dirFiles, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	for _, file := range dirFiles {
		f, err := pinot_crds.CRDs.Open(file.Name())
		if err != nil {
			return err
		}

		err = r.load(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *CRDReconciler) load(f io.Reader) error {
	var b bytes.Buffer

	var yamls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			yamls = append(yamls, b.String())
			b.Reset()
		} else {
			if _, err := b.WriteString(line); err != nil {
				return err
			}
			if _, err := b.WriteString("\n"); err != nil {
				return err
			}
		}
	}
	if s := strings.TrimSpace(b.String()); s != "" {
		yamls = append(yamls, s)
	}

	for _, yaml := range yamls {
		s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme,
			scheme.Scheme)

		obj, _, err := s.Decode([]byte(yaml), nil, nil)
		if err != nil {
			continue
		}

		if crd, ok := obj.(*apiextensionsv1.CustomResourceDefinition); ok {
			crd.Status = apiextensionsv1.CustomResourceDefinitionStatus{}
			crd.SetGroupVersionKind(apiextensionsv1.SchemeGroupVersion.WithKind("CustomResourceDefinition"))
			crd.SetLabels(util.MergeStringMaps(crd.GetLabels(), map[string]string{
				createdByLabel: createdBy,
			}))
			r.crds = append(r.crds, crd)
			continue
		}
	}

	return nil
}

func (r *CRDReconciler) Reconcile(config *pinotv1alpha1.Pinot, log logr.Logger) error {
	log = log.WithValues("component", componentName)

	for _, obj := range r.crds {
		var name, kind string
		if crd, ok := obj.(*apiextensionsv1.CustomResourceDefinition); ok {
			name = crd.Name
			kind = crd.Spec.Names.Kind
		} else if crd, ok := obj.(*apiextensionsv1beta1.CustomResourceDefinition); ok {
			name = crd.Name
			kind = crd.Spec.Names.Kind
		} else {
			log.Error(errors.New("invalid GVK"), "cannot reconcile CRD", "gvk", obj.GetObjectKind().GroupVersionKind())
			continue
		}

		crd := obj.DeepCopyObject()
		current := obj.DeepCopyObject()
		err := k8sutil.SetResourceRevision(crd, r.revision)
		if err != nil {
			return emperror.Wrap(err, "could not set resource revision")
		}
		log := log.WithValues("kind", kind)
		err = r.client.Get(context.Background(), client.ObjectKey{
			Name: name,
		}, current)
		if err != nil && !apierrors.IsNotFound(err) {
			return emperror.WrapWith(err, "getting CRD failed", "kind", kind)
		}
		if apierrors.IsNotFound(err) {
			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(crd); err != nil {
				log.Error(err, "Failed to set last applied annotation", "crd", crd)
			}
			if err := r.client.Create(context.Background(), crd); err != nil {
				return emperror.WrapWith(err, "creating CRD failed", "kind", kind)
			}
			log.Info("CRD created")
		} else {
			if ok, err := k8sutil.CheckResourceRevision(current, fmt.Sprintf("<=%s", r.revision)); !ok {
				if err != nil {
					log.Error(err, "could not check resource revision")
				} else {
					log.V(1).Info("CRD is too new for us")
				}
				continue
			}
			metaAccessor := meta.NewAccessor()
			currentResourceVersion, err := metaAccessor.ResourceVersion(current)
			if err != nil {
				return err
			}

			metaAccessor.SetResourceVersion(crd, currentResourceVersion)

			patchResult, err := patch.DefaultPatchMaker.Calculate(current, crd, patch.IgnoreStatusFields())
			if err != nil {
				log.Error(err, "could not match objects", "kind", kind)
			} else if patchResult.IsEmpty() {
				log.V(1).Info("CRD is in sync")
				continue
			} else {
				log.V(1).Info("resource diffs",
					"patch", string(patchResult.Patch),
					"current", string(patchResult.Current),
					"modified", string(patchResult.Modified),
					"original", string(patchResult.Original))
			}

			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(crd); err != nil {
				log.Error(err, "Failed to set last applied annotation", "crd", crd)
			}

			if err := r.client.Update(context.Background(), crd); err != nil {
				errorMessage := "updating CRD failed, consider updating the CRD manually if needed"
				r.recorder.Eventf(
					config,
					"Warning",
					"PinotCRDUpdateFailure",
					errorMessage,
					"kind",
					kind,
				)
				return emperror.WrapWith(err, errorMessage, "kind", kind)
			}
			log.Info("CRD updated")
		}
	}

	log.Info("Reconciled")

	return nil
}

func GetWatchPredicateForCRDs() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			if e.Meta.GetLabels()[createdByLabel] == createdBy {
				return true
			}
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.MetaOld.GetLabels()[createdByLabel] == createdBy || e.MetaNew.GetLabels()[createdByLabel] == createdBy {
				return true
			}
			return false
		},
	}
}
