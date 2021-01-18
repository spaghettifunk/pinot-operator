package k8sutil

import (
	"github.com/Masterminds/semver/v3"
	"github.com/goph/emperror"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ResourceRevisionLabel = "resource.alpha.apache.io/revision"
)

func SetResourceRevision(obj runtime.Object, revision string) error {
	m, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	labels := m.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	m.SetLabels(util.MergeStringMaps(labels, map[string]string{
		ResourceRevisionLabel: revision,
	}))

	return nil
}

func GetResourceRevision(obj runtime.Object) (string, error) {
	m, err := meta.Accessor(obj)
	if err != nil {
		return "", err
	}

	return m.GetLabels()[ResourceRevisionLabel], nil
}

func CheckResourceRevision(obj runtime.Object, revisionConstraint string) (bool, error) {
	semverConstraint, err := semver.NewConstraint(revisionConstraint)
	if err != nil {
		return false, emperror.Wrap(err, "could not create semver constraint")
	}
	currentRevision, err := GetResourceRevision(obj)
	if err != nil {
		return false, emperror.Wrap(err, "could not get current revision")
	}

	if currentRevision != "" {
		if currentSemver, err := semver.NewVersion(currentRevision); err == nil && !semverConstraint.Check(currentSemver) {
			return false, nil
		}
	}

	return true, nil
}
