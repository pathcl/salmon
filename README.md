salmon
======

### Problem

- You need to retrieve all expiration dates from ingresses using tls on kubernetes

### How?

    $ GO111MODULE=on go run -v github.com/pathcl/salmon
    go: finding github.com/pathcl/salmon latest
    2019/10/15 16:02:13 example.dev.k8s.domain.tld {"cn":"*.dev.k8s.domain.tld","expires":"2019-11-24T13:33:49Z","issuer":"Let's Encrypt Authority X3"}
    2019/10/15 16:02:13 exampletwo.dev.k8s.domain.tld {"cn":"*.dev.k8s.domain.tld","expires":"2019-11-24T13:33:49Z","issuer":"Let's Encrypt Authority X3"}
    2019/10/15 16:02:13 examplethree.dev.k8s.domain.tld {"cn":"*.dev.k8s.domain.tld","expires":"2019-11-24T13:33:49Z","issuer":"Let's Encrypt Authority X3"}