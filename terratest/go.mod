module k8gbterratest

go 1.19

require (
	github.com/AbsaOSS/gopkg v0.1.3
	github.com/gruntwork-io/terratest v0.41.16
	github.com/stretchr/testify v1.8.1
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.20.6
	k8s.io/apimachinery v0.20.6
)

require (
	cloud.google.com/go/compute v1.12.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.1 // indirect
	github.com/aws/aws-sdk-go v1.44.122 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-errors/errors v1.0.2-0.20180813162953-d98b870cc4e0 // indirect
	github.com/go-logr/logr v0.2.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/gruntwork-io/go-commons v0.8.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/lixiangzhong/dnsutil v0.0.0-20191203032812-75ad39d2945a // indirect
	github.com/mattn/go-zglob v0.0.2-0.20190814121620-e3c945676326 // indirect
	github.com/miekg/dns v1.1.31 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/otp v1.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/urfave/cli v1.22.2 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/oauth2 v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/client-go v0.20.6 // indirect
	k8s.io/klog/v2 v2.4.0 // indirect
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.0.3 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	// golang.org/x/crypto
	golang.org/x/crypto v0.0.0-20180904163835-0709b304e793 => golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85 => golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529 => golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5 => golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8 => golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 => golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 => golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
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
	// golang.org/x/text
	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0 => golang.org/x/text v0.3.8
	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c => golang.org/x/text v0.3.8
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.8
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2 => golang.org/x/text v0.3.8
	golang.org/x/text v0.3.2 => golang.org/x/text v0.3.8
	golang.org/x/text v0.3.3 => golang.org/x/text v0.3.8
	golang.org/x/text v0.3.4 => golang.org/x/text v0.3.8
	golang.org/x/text v0.3.5 => golang.org/x/text v0.3.8
	golang.org/x/text v0.3.6 => golang.org/x/text v0.3.8
	golang.org/x/text v0.3.7 => golang.org/x/text v0.3.8
)

exclude (
	github.com/Azure/go-autorest/autorest/adal v0.9.0
	github.com/Azure/go-autorest/autorest/adal v0.9.2
	github.com/Azure/go-autorest/autorest/adal v0.9.5
)
