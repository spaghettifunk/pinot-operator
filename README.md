# Apache Pinot Kuberentes Operator

[![Go Report Card](https://goreportcard.com/badge/github.com/spaghettifunk/pinot-operator)](https://goreportcard.com/report/github.com/spaghettifunk/pinot-operator)

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

The operator is based on the `kubebuilder` project and it has being scaffolded with it.

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
