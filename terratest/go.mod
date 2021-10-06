module k8gbterratest

go 1.16

require (
	github.com/AbsaOSS/gopkg v0.1.2
	github.com/gruntwork-io/terratest v0.37.12
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.6
	k8s.io/apimachinery v0.20.6
)

// CVE-2021-41103 (https://github.com/advisories/GHSA-c2h3-6mxw-7mvq)
// CVE-2020-27813 (https://github.com/advisories/GHSA-3xh2-74w9-5vxm)
// CVE-2020-26160 (https://github.com/advisories/GHSA-w73w-5m7g-f7qc)
require (
	github.com/Azure/go-autorest/autorest/adal v0.9.16 // indirect
	github.com/containerd/containerd v1.5.7 // indirect
	github.com/google/go-containerregistry v0.6.0 // indirect
	github.com/spf13/cobra v1.2.1 // indirect
	go.etcd.io/etcd v3.3.26+incompatible // indirect
)

replace (
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.9.16
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.7
	github.com/google/go-containerregistry v0.0.0-20200110202235-f4fb41bf00a3 => github.com/google/go-containerregistry v0.6.0
	github.com/spf13/cobra v1.0.0 => github.com/spf13/cobra v1.2.1
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489 => go.etcd.io/etcd v3.3.26+incompatible
)

exclude  (
	github.com/google/go-containerregistry v0.0.0-20200110202235-f4fb41bf00a3
	github.com/containerd/containerd v1.5.0-beta.1
	github.com/containerd/containerd v1.5.0-beta.3
	github.com/containerd/containerd v1.5.0-beta.4
	github.com/containerd/containerd v1.5.0-rc.0
	github.com/containerd/containerd v1.5.1
	github.com/containerd/containerd v1.5.2
	github.com/Azure/go-autorest/autorest/adal v0.9.5
	github.com/Azure/go-autorest/autorest/adal v0.9.2
	github.com/Azure/go-autorest/autorest/adal v0.9.0
	github.com/Azure/go-autorest/autorest/adal v0.8.2
	github.com/Azure/go-autorest/autorest/adal v0.5.0
	github.com/spf13/cobra v1.0.0
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489
)
