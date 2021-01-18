# runs the target list by default
.DEFAULT_GOAL = list

.PHONY: list

# Image URL to use all building/pushing image targets
OPERATOR_IMAGE = davideberdin/pinot-operator:v0.0.1
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS = crd:crdVersions={v1},maxDescLen=0,trivialVersions=false

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Insert a comment starting with '##' after a target, and it will be printed by 'make' and 'make list'
list:    ## list Makefile targets
	@echo "The most used targets: \n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

unit-tests: install-tools generate fmt vet manifests ## Run unit tests
	ginkgo -r --randomizeAllSpecs api/ internal/	

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

debug: generate fmt vet manifests

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${OPERATOR_IMAGE}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=operator-role paths="./api/...;./controllers/..." output:crd:artifacts:config=config/crd/bases	

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: api-reference
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${OPERATOR_IMAGE}

# Push the docker image
docker-push:
	docker push ${OPERATOR_IMAGE}

git-commit-sha:
ifeq ("", git diff --stat)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
else
GIT_COMMIT=$(shell git rev-parse --short HEAD)-
endif

docker-build-dev: check-env-docker-repo  git-commit-sha
	docker build --build-arg=GIT_COMMIT=$(GIT_COMMIT) -t $(DOCKER_REGISTRY_SERVER)/$(OPERATOR_IMAGE):$(GIT_COMMIT) .
	docker push $(DOCKER_REGISTRY_SERVER)/$(OPERATOR_IMAGE):$(GIT_COMMIT)

check-env-docker-repo: check-env-registry-server
ifndef OPERATOR_IMAGE
	$(error OPERATOR_IMAGE is undefined: path to the Operator image within the registry specified in DOCKER_REGISTRY_SERVER (e.g. davideberdin/pinot-operator - without leading slash))
endif

check-env-docker-credentials: check-env-registry-server
ifndef DOCKER_REGISTRY_USERNAME
	$(error DOCKER_REGISTRY_USERNAME is undefined: Username for accessing the docker registry)
endif
ifndef DOCKER_REGISTRY_PASSWORD
	$(error DOCKER_REGISTRY_PASSWORD is undefined: Password for accessing the docker registry)
endif
ifndef DOCKER_REGISTRY_SECRET
	$(error DOCKER_REGISTRY_SECRET is undefined: Name of Kubernetes secret in which to store the Docker registry username and password)
endif

check-env-registry-server:
ifndef DOCKER_REGISTRY_SERVER
	$(error DOCKER_REGISTRY_SERVER is undefined: URL of docker registry containing the Operator image (e.g. registry.my-company.com))
endif

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

install-tools:
	go mod download
	grep _ tools/tools.go | awk -F '"' '{print $$2}' | xargs -t go install

# documentation
api-reference: install-tools ## Generate API reference documentation
	crd-ref-docs \
		--source-path ./api/v1alpha1 \
		--config ./docs/api/autogen/config.yaml \
		--templates-dir ./docs/api/autogen/templates \
		--output-path ./docs/api/pinot.apache.io.ref.asciidoc \
		--max-depth 30

# Builds a single-file installation manifest to deploy the Operator
generate-installation-manifest:
	mkdir -p releases
	kustomize build config/installation/ > releases/pinot-cluster-operator.yaml