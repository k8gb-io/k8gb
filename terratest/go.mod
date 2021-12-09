module k8gbterratest

go 1.16

require (
	github.com/AbsaOSS/gopkg v0.1.2
	github.com/gruntwork-io/terratest v0.38.6
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.6
	k8s.io/apimachinery v0.20.6
)

replace (
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.9.16
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.8
	github.com/opencontainers/image-spec v1.0.1 => github.com/opencontainers/image-spec v1.0.2
	github.com/opencontainers/runc v1.0.2 => github.com/opencontainers/runc v1.0.3
)

exclude (
	github.com/Azure/go-autorest/autorest/adal v0.5.0
	github.com/Azure/go-autorest/autorest/adal v0.8.2
	github.com/Azure/go-autorest/autorest/adal v0.9.0
	github.com/Azure/go-autorest/autorest/adal v0.9.2
	github.com/Azure/go-autorest/autorest/adal v0.9.5
	github.com/containerd/containerd v1.4.9
	github.com/containerd/containerd v1.5.0-beta.1
	github.com/containerd/containerd v1.5.0-beta.3
	github.com/containerd/containerd v1.5.0-beta.4
	github.com/containerd/containerd v1.5.0-rc.0
	github.com/containerd/containerd v1.5.1
	github.com/containerd/containerd v1.5.2
)
