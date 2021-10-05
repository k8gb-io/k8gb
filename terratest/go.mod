module k8gbterratest

go 1.16

require (
	github.com/AbsaOSS/gopkg v0.1.2
	github.com/gruntwork-io/terratest v0.32.17
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5 // indirect
)

replace github.com/containerd/containerd v1.3.0 => github.com/containerd/containerd v1.4.3
