module github.com/spaghettifunk/pinot-operator

go 1.15

require (
	emperror.dev/errors v0.8.0
	github.com/Masterminds/semver/v3 v3.1.0
	github.com/banzaicloud/k8s-objectmatcher v1.4.0
	github.com/elastic/crd-ref-docs v0.0.6
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/strfmt v0.20.0
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/goph/emperror v0.17.2
	github.com/hoisie/mustache v0.0.0-20160804235033-6375acf62c69
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200627165143-92b8a710ab6c
	github.com/spaghettifunk/pinot-go-client v0.0.0-20210120142358-a99a0c96c73c
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools v0.0.0-20200616195046-dc31b401abb5 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/controller-tools v0.2.5 // indirect
)
