# Apache Pinot Kubernetes Operator

A kubernetes operator for Apache Pinot

## Install kubebuilder

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
