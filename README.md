# Apache Pinot Kuberentes Operator

[![Go Report Card](https://goreportcard.com/badge/github.com/spaghettifunk/pinot-operator)](https://goreportcard.com/report/github.com/spaghettifunk/pinot-operator)
[![license](https://img.shields.io/github/license/apache/pinot.svg)](LICENSE)

**Project status: _alpha_** Not all planned features are completed. The API, spec, status and other user facing objects may change, and not in a backward compatible way.

## Overview

The Pinot Operator provides [Kubernetes](https://kubernetes.io/) native deployment and management of
[Apache Prometheus]() and related components. The purpose of this project is to
simplify and automate the configuration of a Apache Pinot stack for Kubernetes clusters.

The Pinot operator includes, but is not limited to, the following features:

- **Kubernetes Custom Resources**: Use Kubernetes custom resources to deploy and manage Apache Pnot and related components.

- **Simplified Deployment Configuration**: Configure the fundamentals of Brokers/Controller/Server like versions, persistence,
  retention policies, and replicas from a native Kubernetes resource.

The operator has been largely inspired by the [BanzaiCloud istio-operator](https://github.com/banzaicloud/istio-operator) and the [RabbitMQ cluster-operator](https://github.com/rabbitmq/cluster-operator). They are great resources to learn how to create operators.

## Quickstart

If you have a running Kubernetes cluster and `kubectl` configured to access it, run the following command to install the operator:

```bash
kubectl apply -f https://github.com/spaghettifunk/pinot-operator/releases/latest/download/pinot-cluster-operator.yaml
```

Then you can deploy a Pinot cluster:

```bash
kubectl apply -f https://raw.githubusercontent.com/spaghettifunk/pinot-operator/main/docs/examples/hello-world/pinot.yaml
```

[![asciicast](https://asciinema.org/a/385228.svg)](https://asciinema.org/a/385228)

## How to develop

The operator is based on the `kubebuilder` project and it has being scaffolded with it. To make it run, you need to do a few steps:

1. `make generate` to generate the `deepcopy` files
2. `make manifests` to generate the correct CRDs
3. Initiate KinD with `kind create cluster`. If you do not have KinD, check this [page](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) to install it
4. `make install` to deploy the CRDs to your cluster
5. `WATCH_NAMESPACE=pinot-system POD_NAMESPACE=pinot-system make run` to run the pinot-controller locally

If you want to stop the controller, press `CTRL-C` and wait 30 seconds for the `stop handler` to complete.

### Install kubebuilder

If you are on `mac os x` you can install `kubebuilder` with `brew install kubebuilder`. Bare in mind that the `brew` command install only the binary of tool.
To run the tests, you need some extra packages within your `$PATH`. For convenience, there is a file called `hack/install_kubebuilder_pkg.sh` that will pull the extra files and put it
in the right directory for you. Before you can actually run the tests, you **must** have them in your machine.

## Testing

Run `make test` to launch the test suite. If you see an error similar to this one

```
Failure [0.007 seconds]
[BeforeSuite] BeforeSuite
/Users/davideberdin/go/src/github.com/spaghettifunk/pinot-operator/controllers/suite_test.go:52
  Unexpected error:
      <*fmt.wrapError | 0xc00038d9c0>: {
          msg: "failed to start the controlplane. retried 5 times: fork/exec /usr/local/kubebuilder/bin/etcd: no such file or directory",
          err: {
              Op: "fork/exec",
              Path: "/usr/local/kubebuilder/bin/etcd",
              Err: 0x2,
          },
      }
      failed to start the controlplane. retried 5 times: fork/exec /usr/local/kubebuilder/bin/etcd: no such file or directory
  occurred
  /Users/davideberdin/go/src/github.com/spaghettifunk/pinot-operator/controllers/suite_test.go:62
------------------------------
```

you need to run the script called `hack/install_kubebuilder_pkg.sh`.

## Versioning

Apache Pinot Kubernetes Operator follows non-strict [semver](https://semver.org/).

[The versioning guidelines document](version_guidelines.md) contains guidelines
on how we implement non-strict semver. The version number MAY or MAY NOT follow the semver rules. Hence, we highly recommend to read
the release notes to understand the changes and their potential impact for any release.

## Contributing

This project follows the typical GitHub pull request model. Before starting any work, please either comment on an [existing issue](https://github.com/spaghettifunk/pinot-operator/issues), or file a new one.

Please read [contribution guidelines](CONTRIBUTING.md) if you are interested in contributing to this project.
