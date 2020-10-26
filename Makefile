REPO ?= absaoss/k8gb
# Current Operator version
VERSION ?= $$(helm show chart chart/k8gb/|awk '/appVersion:/ {print $$2}')
K8GB_IMAGE_TAG  ?= v$(VERSION)
# Default bundle image tag
BUNDLE_IMG ?= controller-bundle:$(VERSION)
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Image URL to use all building/pushing image targets
IMG ?= $(REPO):$(K8GB_IMAGE_TAG)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run tests
test: lint generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# Build and push the docker image exclusively for testing using commit hash
docker-test-build-push: test
	$(call docker-test-build-push)

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

# Generate bundle manifests and metadata, then validate generated files.
.PHONY: bundle
bundle: manifests
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

## Special k8gb make part

VALUES_YAML ?= chart/k8gb/values.yaml
HELM_ARGS ?=
ETCD_DEBUG_IMAGE ?= quay.io/coreos/etcd:v3.2.25
GSLB_DOMAIN ?= cloud.example.com
HOST_ALIAS_IP1 ?= 172.17.0.9
HOST_ALIAS_IP2 ?= 172.17.0.5
K8GB_IMAGE_REPO ?= absaoss/k8gb
K8GB_IMAGE_TAG  ?= v$(VERSION)
K8GB_COREDNS_IP ?= kubectl get svc k8gb-coredns -n k8gb -o custom-columns='IP:spec.clusterIP' --no-headers
PODINFO_IMAGE_REPO ?= stefanprodan/podinfo
COMMIT_HASH ?= $(shell git rev-parse --short HEAD)

.PHONY: debug-local
debug-local: create-test-ns
	kubectl apply -f ./deploy/crds/k8gb.absa.oss_gslbs_crd.yaml
	kubectl apply -f ./deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml
	operator-sdk run --local --namespace=test-gslb --enable-delve

.PHONY: lint
lint:
	staticcheck ./...
	errcheck ./...
	golint '-set_exit_status=1' ./...

.PHONY: terratest
terratest:
	cd terratest/test/ && go mod download && go test -v

.PHONY: dns-tools
dns-tools:
	kubectl -n k8gb get svc k8gb-coredns
	kubectl -n k8gb run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools

.PHONY: dns-smoke-test
dns-smoke-test:
	kubectl -n k8gb run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools --command -- /usr/bin/dig @k8gb-coredns app3.cloud.example.com

.PHONY: deploy-local-cluster
deploy-local-cluster:
	kind create cluster --config=deploy/kind/cluster.yaml --name test-gslb1

.PHONY: deploy-two-local-clusters
deploy-two-local-clusters:
	kind create cluster --config=deploy/kind/cluster.yaml --name test-gslb1
	kind create cluster --config=deploy/kind/cluster2.yaml --name test-gslb2

.PHONY: use-first-context
use-first-context:
	kubectl config use-context kind-test-gslb1

.PHONY: use-second-context
use-second-context:
	kubectl config use-context kind-test-gslb2

.PHONY: deploy-first-k8gb
deploy-first-k8gb: HELM_ARGS = --set k8gb.hostAlias.enabled=true --set k8gb.hostAlias.ip="$(HOST_ALIAS_IP1)" --set k8gb.imageRepo=$(K8GB_IMAGE_REPO)
deploy-first-k8gb: deploy-gslb-operator deploy-local-ingress

.PHONY: deploy-second-k8gb
deploy-second-k8gb: HELM_ARGS = --set k8gb.hostAlias.enabled=true --set k8gb.clusterGeoTag="us" --set k8gb.extGslbClustersGeoTags="eu" --set k8gb.hostAlias.hostnames="{test-gslb-ns-eu.example.com,test-gslb-failover-ns-eu.example.com}" --set k8gb.hostAlias.ip="$(HOST_ALIAS_IP2)" --set k8gb.imageRepo=$(K8GB_IMAGE_REPO)
deploy-second-k8gb: deploy-gslb-operator deploy-local-ingress

.PHONY: deploy-full-local-setup
deploy-full-local-setup: deploy-two-local-clusters
	ADDITIONAL_TARGETS=deploy-test-apps VERSION=$(VERSION) ./deploy/full.sh

.PHONY: destroy-full-local-setup
destroy-full-local-setup: destroy-two-local-clusters

.PHONY: destroy-local-cluster
destroy-local-cluster:
	kind delete cluster --name test-gslb1

.PHONY: destroy-two-local-clusters
destroy-two-local-clusters:
	kind delete cluster --name test-gslb1
	kind delete cluster --name test-gslb2

.PHONY: create-k8gb-ns
create-k8gb-ns:
	kubectl apply -f deploy/namespace.yaml

.PHONY: create-test-ns
create-test-ns:
	kubectl apply -f deploy/crds/test-namespace.yaml

.PHONY: deploy-local-ingress
deploy-local-ingress: create-k8gb-ns
	helm repo add --force-update stable https://kubernetes-charts.storage.googleapis.com
	helm repo update
	helm -n k8gb upgrade -i nginx-ingress stable/nginx-ingress --version 1.41.1 -f deploy/ingress/nginx-ingress-values.yaml

.PHONY: wait-for-nginx-ingress-ready
wait-for-nginx-ingress-ready:
	kubectl -n k8gb wait --for=condition=Ready pod -l app=nginx-ingress --timeout=600s

.PHONY: deploy-gslb-operator
deploy-gslb-operator: create-k8gb-ns
	cd chart/k8gb && helm dependency update
	helm -n k8gb upgrade -i k8gb chart/k8gb -f $(VALUES_YAML) $(HELM_ARGS)

.PHONY: wait-for-gslb-ready
wait-for-gslb-ready:
	kubectl -n k8gb wait --for=condition=Ready pod -l app=etcd --timeout=600s

# workaround until https://github.com/crossplaneio/crossplane/issues/1170 solved
.PHONY: deploy-gslb-operator-14
deploy-gslb-operator-14: create-k8gb-ns
	cd chart/k8gb && helm dependency update
	helm -n k8gb template k8gb chart/k8gb -f $(VALUES_YAML) | kubectl -n k8gb --validate=false apply -f -

.PHONY: deploy-gslb-cr
deploy-gslb-cr: create-test-ns
	$(call apply-cr,deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml)
	$(call apply-cr,deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr_failover.yaml)

.PHONY: deploy-test-apps
deploy-test-apps: create-test-ns
	kubectl apply -f deploy/test-apps
	helm repo add podinfo https://stefanprodan.github.io/podinfo
	helm upgrade --install frontend --namespace test-gslb -f deploy/test-apps/podinfo/podinfo-values.yaml --set ui.message="`$(call get-cluster-geo-tag)`" --set image.repository="$(PODINFO_IMAGE_REPO)" podinfo/podinfo

.PHONY: clean-test-apps
clean-test-apps:
	kubectl delete -f deploy/test-apps
	helm -n test-gslb uninstall backend
	helm -n test-gslb uninstall frontend

.PHONY: debug-test-etcd
debug-test-etcd:
	kubectl run --rm -i --tty --env="ETCDCTL_API=3" --env="ETCDCTL_ENDPOINTS=http://etcd-cluster-client:2379" --namespace k8gb etcd-test --image "$(ETCD_DEBUG_IMAGE)" --restart=Never -- /bin/sh

.PHONY: infoblox-secret
infoblox-secret:
	kubectl -n k8gb create secret generic external-dns \
	    --from-literal=EXTERNAL_DNS_INFOBLOX_WAPI_USERNAME=$${WAPI_USERNAME} \
	    --from-literal=EXTERNAL_DNS_INFOBLOX_WAPI_PASSWORD=$${WAPI_PASSWORD}

.PHONY: init-failover
init-failover:
	$(call init-test-strategy, "deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr_failover.yaml")

.PHONY: init-round-robin
init-round-robin:
	$(call init-test-strategy, "deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml")

.PHONY: stop-test-app
stop-test-app:
	$(call testapp-set-replicas,0)

.PHONY: start-test-app
start-test-app:
	$(call testapp-set-replicas,2)

.PHONY: test-round-robin
test-round-robin:
	@$(call hit-testapp-host, "app3.cloud.example.com")

.PHONY: test-failover
test-failover:
	@$(call hit-testapp-host, "failover.cloud.example.com")

.PHONY: demo-roundrobin
demo-roundrobin:
	@$(call demo-host, "app3.cloud.example.com")

.PHONY: demo-failover
demo-failover:
	@$(call demo-host, "failover.cloud.example.com")

.PHONY: version
version:
	@echo $(VERSION)

define testapp-set-replicas
	kubectl scale deployment frontend-podinfo -n test-gslb --replicas=$1
endef

define hit-testapp-host
	kubectl run -it --rm busybox --restart=Never --image=busybox -- sh -c \
	"echo 'nameserver `$(K8GB_COREDNS_IP)`' > /etc/resolv.conf && \
	wget -qO - $1"
endef

define demo-host
	kubectl run -it --rm k8gbdemo --restart=Never --image=absaoss/k8gbdemocurl  \
	"`$(K8GB_COREDNS_IP)`" $1
endef

define init-test-strategy
 	kubectl config use-context kind-test-gslb2
 	kubectl apply -f $1
 	kubectl config use-context kind-test-gslb1
 	kubectl apply -f $1
 	$(call testapp-set-replicas,2)
endef

define get-cluster-geo-tag
	kubectl -n k8gb describe deploy k8gb |  awk '/CLUSTER_GEO_TAG/ { printf $$2 }'
endef

define apply-cr
	sed -i 's/cloud\.example\.com/$(GSLB_DOMAIN)/g' "$1"
	kubectl apply -f "$1"
	git checkout -- "$1"
endef

define docker-test-build-push
    docker build . -t k8gb:$(COMMIT_HASH)
    docker tag k8gb:$(COMMIT_HASH) $(REPO):v$(COMMIT_HASH)
    docker push $(REPO):v$(COMMIT_HASH)
    sed -i "s/$(VERSION)/$(COMMIT_HASH)/g" chart/k8gb/Chart.yaml
endef
