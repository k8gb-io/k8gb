module github.com/AbsaOSS/k8gb

go 1.16

require (
	github.com/AbsaOSS/gopkg v0.1.2
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.4.0
	github.com/golang/mock v1.5.0
	github.com/infobloxopen/infoblox-go-client v1.1.0
	github.com/miekg/dns v1.1.42
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.21.0
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/external-dns v0.8.0
)

replace golang.org/x/crypto => golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
