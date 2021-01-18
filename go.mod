module github.com/spaghettifunk/pinot-operator

go 1.15

require (
	github.com/banzaicloud/istio-operator v0.0.0-20210118101045-2a25c679bca7
	github.com/banzaicloud/k8s-objectmatcher v1.4.0
	github.com/elastic/crd-ref-docs v0.0.6
	github.com/go-logr/logr v0.1.0
	github.com/goph/emperror v0.17.2
	github.com/hoisie/mustache v0.0.0-20160804235033-6375acf62c69
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/controller-tools v0.2.5 // indirect
)
