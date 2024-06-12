module github.com/k8gb-io/k8gb

go 1.22.3

require (
	github.com/AbsaOSS/env-binder v1.0.1
	github.com/AbsaOSS/gopkg v0.1.3
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v1.2.3
	github.com/infobloxopen/infoblox-go-client v1.1.1
	github.com/miekg/dns v1.1.52
	github.com/prometheus/client_golang v1.14.0
	github.com/rs/zerolog v1.21.0
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.7.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	go.uber.org/mock v0.4.0
	k8s.io/api v0.26.3
	k8s.io/apimachinery v0.26.3
	k8s.io/client-go v0.26.3
	sigs.k8s.io/controller-runtime v0.14.5
	sigs.k8s.io/external-dns v0.13.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.10.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20220622183110-fd043fe589d2 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.3.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220804142021-4e6b2dfa6612 // indirect
	google.golang.org/grpc v1.49.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.26.1 // indirect
	k8s.io/component-base v0.26.1 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	// golang.org/x/net
	golang.org/x/net v0.0.0-20180724234803-3673e40ba225 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20180826012351-8a410e7b638d => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20180906233101-161cd47e91fd => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190108225652-1e06a53dbb7e => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190501004415-9ce7a6920f09 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200222125558-5a598a2470a0 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200501053045-e0ff5e5a1de5 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f => golang.org/x/net v0.7.0
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b => golang.org/x/net v0.7.0
)
