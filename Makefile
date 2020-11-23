###############################
#		CONSTANTS
###############################
CLUSTER_GSLB1 = test-gslb1
CLUSTER_GSLB2 = test-gslb2
GSLB_DOMAIN ?= cloud.example.com
REPO = absaoss/k8gb
VALUES_YAML ?= chart/k8gb/values.yaml
PODINFO_IMAGE_REPO ?= stefanprodan/podinfo
HELM_ARGS ?= --set k8gb.clusterGeoTag='us' --set k8gb.extGslbClustersGeoTags='eu' --set k8gb.hostAlias.hostnames='{test-gslb-ns-eu.example.com,test-gslb-failover-ns-eu.example.com}'
K8GB_COREDNS_IP ?= kubectl get svc k8gb-coredns -n k8gb -o custom-columns='IP:spec.clusterIP' --no-headers
ETCD_DEBUG_IMAGE ?= quay.io/coreos/etcd:v3.2.25

# terratest
GITACTION_TERRATEST_DOCKER_REPO_PORT ?= 5000
GITACTION_TERRATEST_DOCKER_REPO_NAME ?= kind-registry
GITACTION_TERRATEST_DOCKER_REPO_IMAGE ?= registry:2

ifndef NO_COLOR
YELLOW=\033[0;33m
CYAN=\033[1;36m
# no color
NC=\033[0m
endif

NO_VALUE ?= no_value

###############################
#		VARIABLES
###############################

# shell script inspects whether kind registry is running or not
KIND_REGISTRY_RUNNING ?= $(shell docker inspect -f '{{.State.Running}}' "test-gslb2-control-plane" || true)

VERSION ?= $(shell helm show chart chart/k8gb/|awk '/appVersion:/ {print $$2}')

# image URL to use all building/pushing image targets
IMG ?= $(REPO):v$(VERSION)

# default bundle image tag
BUNDLE_IMG ?= controller-bundle:$(VERSION)

# options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# create GOBIN if not specified
ifndef GOBIN
GOBIN=$(shell go env GOPATH)/bin
endif

CONTROLLER_GEN_PATH ?= $(shell which controller-gen || echo $(NO_VALUE))

KUSTOMIZE_PATH ?= $(shell which kustomize || echo $(NO_VALUE))

COMMIT_HASH ?= $(shell git rev-parse --short HEAD)

###############################
#		TARGETS
###############################

all: manager

.PHONY: clean-test-apps
clean-test-apps:
	kubectl delete -f deploy/test-apps
	helm -n test-gslb uninstall backend
	helm -n test-gslb uninstall frontend

# see: https://dev4devs.com/2019/05/04/operator-framework-how-to-debug-golang-operator-projects/
.PHONY: debug-idea
debug-idea: export WATCH_NAMESPACE=test-gslb
debug-idea:
	$(call debug,debug --headless --listen=:2345 --api-version=2)

.PHONY: debug-test-etcd
debug-test-etcd:
	kubectl run --rm -i --tty --env="ETCDCTL_API=3" --env="ETCDCTL_ENDPOINTS=http://etcd-cluster-client:2379" --namespace k8gb etcd-test --image "$(ETCD_DEBUG_IMAGE)" --restart=Never -- /bin/sh

.PHONY: demo-roundrobin
demo-roundrobin:
	@$(call demo-host, "app3.cloud.example.com")

.PHONY: demo-failover
demo-failover:
	@$(call demo-host, "failover.cloud.example.com")

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
.PHONY: deploy
deploy:
	$(call manifest)
	cd config/manager && $(KUSTOMIZE_PATH) edit set image controller=$(IMG)
	$(KUSTOMIZE_PATH) build config/default | kubectl apply -f -

# spin-up local environment
.PHONY: deploy-full-local-setup
deploy-full-local-setup:
	$(call create-local-cluster,$(CLUSTER_GSLB1),"deploy/kind/cluster.yaml")
	$(call create-local-cluster,$(CLUSTER_GSLB2),"deploy/kind/cluster2.yaml")

	$(call deploy-local-cluster,$(CLUSTER_GSLB1),$(CLUSTER_GSLB2),worker,absaoss/k8gb,)

	$(call deploy-local-cluster,$(CLUSTER_GSLB2),$(CLUSTER_GSLB1),worker,absaoss/k8gb,$(HELM_ARGS))

# triggered by terraform GitHub Action. Clusters already exists. GO is not installed yet
.PHONY: deploy-full-terratest-setup
deploy-full-terratest-setup:
	@echo "\n$(YELLOW)create docker registry$(NC) see: https://kind.sigs.k8s.io/docs/user/local-registry/"
	docker run -d --restart=always -p "$(GITACTION_TERRATEST_DOCKER_REPO_PORT):5000" --name "$(GITACTION_TERRATEST_DOCKER_REPO_NAME)" "$(GITACTION_TERRATEST_DOCKER_REPO_IMAGE)"

	@echo "\n$(YELLOW)build k8gb docker and push to registry $(NC)"
	docker build . -t k8gb:$(COMMIT_HASH)
	docker tag k8gb:$(COMMIT_HASH) localhost:$(GITACTION_TERRATEST_DOCKER_REPO_PORT)/k8gb:v$(COMMIT_HASH)
	docker push localhost:$(GITACTION_TERRATEST_DOCKER_REPO_PORT)/k8gb:v$(COMMIT_HASH)

	@echo "\n$(YELLOW)Change version in Chart.yaml $(CYAN) $(VERSION) to $(COMMIT_HASH)$(NC)"
	sed -i "s/$(VERSION)/$(COMMIT_HASH)/g" chart/k8gb/Chart.yaml

	$(call deploy-local-cluster,$(CLUSTER_GSLB1),$(CLUSTER_GSLB2),control-plane,localhost:$(GITACTION_TERRATEST_DOCKER_REPO_PORT)/k8gb,)
	$(call deploy-local-cluster,$(CLUSTER_GSLB2),$(CLUSTER_GSLB1),control-plane,localhost:$(GITACTION_TERRATEST_DOCKER_REPO_PORT)/k8gb,$(HELM_ARGS))

	@echo "\n$(YELLOW)Local cluster $(CYAN)$(CLUSTER_GSLB2) $(NC)"
	kubectl get pods -A
	@echo "\n$(YELLOW)Local cluster $(CYAN)$(CLUSTER_GSLB1) $(NC)"
	kubectl config use-context kind-$(CLUSTER_GSLB1)
	kubectl get pods -A

.PHONY: deploy-gslb-operator
deploy-gslb-operator:
	cd chart/k8gb && helm dependency update
	helm -n k8gb upgrade -i k8gb chart/k8gb -f $(VALUES_YAML) $(HELM_ARGS)

# workaround until https://github.com/crossplaneio/crossplane/issues/1170 solved
.PHONY: deploy-gslb-operator-14
deploy-gslb-operator-14:
	kubectl apply -f deploy/namespace.yaml
	cd chart/k8gb && helm dependency update
	helm -n k8gb template k8gb chart/k8gb -f $(VALUES_YAML) | kubectl -n k8gb --validate=false apply -f -

# destroy local test environment
.PHONY: destroy-full-local-setup
destroy-full-local-setup:
	kind delete cluster --name $(CLUSTER_GSLB1)
	kind delete cluster --name $(CLUSTER_GSLB2)

.PHONY: dns-tools
dns-tools:
	@kubectl -n k8gb get svc k8gb-coredns
	@kubectl -n k8gb run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools

.PHONY: dns-smoke-test
dns-smoke-test:
	kubectl -n k8gb run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools --command -- /usr/bin/dig @k8gb-coredns app3.cloud.example.com

# build the docker image
.PHONY: docker-build
docker-build: test
	docker build . -t $(IMG)

# push the docker image
.PHONY: docker-push
docker-push:
	docker push $(IMG)

# build and push the docker image exclusively for testing using commit hash
.PHONY: docker-test-build-push
docker-test-build-push: test
	$(call docker-test-build-push)

.PHONY: init-failover
init-failover:
	$(call init-test-strategy, "deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr_failover.yaml")

.PHONY: init-round-robin
init-round-robin:
	$(call init-test-strategy, "deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml")

# creates infoblox secret in current cluster
.PHONY: infoblox-secret
infoblox-secret:
	kubectl -n k8gb create secret generic external-dns \
		--from-literal=EXTERNAL_DNS_INFOBLOX_WAPI_USERNAME=$${WAPI_USERNAME} \
		--from-literal=EXTERNAL_DNS_INFOBLOX_WAPI_PASSWORD=$${WAPI_PASSWORD}

# install CRDs into a cluster
.PHONY: install
install:
	$(call manifest)
	$(KUSTOMIZE_PATH) build config/crd | kubectl apply -f -

# run all linters from .golangci.yaml; see: https://golangci-lint.run/usage/install/#local-installation
.PHONY: lint
lint:
	golangci-lint run

# retrieves all targets
.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

# build manager binary
.PHONY: manager
manager: lint
	$(call generate)
	go build -o bin/manager main.go

# remove clusters and redeploy
.PHONY: reset
reset:	destroy-full-local-setup deploy-full-local-setup

# run against the configured Kubernetes cluster in ~/.kube/config
.PHONY: run
run: lint
	$(call generate)
	$(call manifest)
	go run ./main.go

.PHONY: stop-test-app
stop-test-app:
	$(call testapp-set-replicas,0)

.PHONY: start-test-app
start-test-app:
	$(call testapp-set-replicas,2)

# run tests
.PHONY: test
test: lint
	$(call generate)
	$(call manifest)
	go test ./... -coverprofile cover.out

.PHONY: test-round-robin
test-round-robin:
	@$(call hit-testapp-host, "app3.cloud.example.com")

.PHONY: test-failover
test-failover:
	@$(call hit-testapp-host, "failover.cloud.example.com")

# executes terra-tests
.PHONY: terratest
terratest:
	cd terratest/test/ && go mod download && go test -v

# uninstall CRDs from a cluster
.PHONY: uninstall
uninstall:
	$(call manifest)
	$(call install-kustomize-if-not-exists)
	$(KUSTOMIZE_PATH) build config/crd | kubectl delete -f -

.PHONY: version
version:
	@echo $(VERSION)

###############################
#		FUNCTIONS
###############################

define create-local-cluster
	@echo "\n$(YELLOW)Deploy local cluster $(CYAN)$1 $(NC)"
	kind create cluster --name=$1 --config=$2
endef

define deploy-local-cluster
	@echo "\n$(YELLOW)Local cluster $(CYAN)$1 $(NC)"
	kubectl config use-context kind-$1

	@echo "\n$(YELLOW)Create namespace $(NC)"
	kubectl apply -f deploy/namespace.yaml

	@echo "\n$(YELLOW)Deploy GSLB operator from $4 $(NC)"
	cd chart/k8gb && helm dependency update
	helm -n k8gb upgrade -i k8gb chart/k8gb -f $(VALUES_YAML) \
		--set k8gb.hostAlias.enabled=true \
		--set k8gb.hostAlias.ip="`$(call get-host-alias-ip,$1,$2,$3)`" \
		--set k8gb.imageRepo=$4 $5

	@echo "\n$(YELLOW)Deploy Ingress $(NC)"
	helm repo add --force-update stable https://charts.helm.sh/stable
	helm repo update
	helm -n k8gb upgrade -i nginx-ingress stable/nginx-ingress \
		--version 1.41.1 -f deploy/ingress/nginx-ingress-values.yaml

	@echo "\n$(YELLOW)Deploy GSLB cr $(NC)"
	kubectl apply -f deploy/crds/test-namespace.yaml
	$(call apply-cr,deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml)
	$(call apply-cr,deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr_failover.yaml)

	@echo "\n$(YELLOW)Deploy test apps $(NC)"
	$(call deploy-test-apps)

	@echo "\n$(YELLOW)Wait for GSLB, Ingress are ready $(NC)"
	$(call wait)

	@echo "\n$(CYAN)$1 $(YELLOW)deployed! $(NC)"
endef

define apply-cr
	sed -i 's/cloud\.example\.com/$(GSLB_DOMAIN)/g' "$1"
	kubectl apply -f "$1"
	git checkout -- "$1"
endef

define deploy-test-apps
	kubectl apply -f deploy/test-apps
	helm repo add podinfo https://stefanprodan.github.io/podinfo
	helm upgrade --install frontend --namespace test-gslb -f deploy/test-apps/podinfo/podinfo-values.yaml \
		--set ui.message="`$(call get-cluster-geo-tag)`" \
		--set image.repository="$(PODINFO_IMAGE_REPO)" \
		podinfo/podinfo
endef

define get-cluster-geo-tag
	kubectl -n k8gb describe deploy k8gb |  awk '/CLUSTER_GEO_TAG/ { printf $$2 }'
endef

# get-host-alias-ip switch to second context ($2), search for IP and switch back to first context ($1)
# function returns one IP address
define get-host-alias-ip
	kubectl config use-context kind-$2 > /dev/null && \
	kubectl get nodes $2-$3 -o custom-columns='IP:status.addresses[0].address' --no-headers && \
	kubectl config use-context kind-$1 > /dev/null
endef

define hit-testapp-host
	kubectl run -it --rm busybox --restart=Never --image=busybox -- sh -c \
	"echo 'nameserver `$(K8GB_COREDNS_IP)`' > /etc/resolv.conf && \
	wget -qO - $1"
endef

define init-test-strategy
 	kubectl config use-context kind-test-gslb2
 	kubectl apply -f $1
 	kubectl config use-context kind-test-gslb1
 	kubectl apply -f $1
 	$(call testapp-set-replicas,2)
endef

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

# waits for NGINX, GSLB are ready
define wait
	kubectl -n k8gb wait --for=condition=Ready pod -l app=nginx-ingress --timeout=600s
	kubectl -n k8gb wait --for=condition=Ready pod -l app=etcd --timeout=600s
endef

define generate
	$(call controller-gen,object:headerFile="hack/boilerplate.go.txt" paths="./...")
endef

define manifest
	$(call controller-gen,crd:trivialVersions=true rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases)
endef

# function retrieves controller-gen path or installs controller-gen@v3.0.0 and retrieve new path in case it is not installed
define controller-gen
	@$(if $(filter $(CONTROLLER_GEN_PATH),$(NO_VALUE)),$(call install-controller-gen),)
	$(CONTROLLER_GEN_PATH) $1
endef

define install-controller-gen
	GO111MODULE=on go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0
	$(eval CONTROLLER_GEN_PATH = $(GOBIN)/controller-gen)
endef

# installs kustomize and sets KUSTOMIZE_PATH if is not specified
define install-kustomize-if-not-exists
	@$(if $(filter $(KUSTOMIZE_PATH),$(NO_VALUE)),$(call install-kustomize),)
endef

define install-kustomize
	GO111MODULE=on go get sigs.k8s.io/kustomize/kustomize/v3@v3.8.6
	$(eval KUSTOMIZE_PATH = $(GOBIN)/kustomize)
endef

define docker-test-build-push
	docker build . -t k8gb:$(COMMIT_HASH)
	docker tag k8gb:$(COMMIT_HASH) $(REPO):v$(COMMIT_HASH)
	docker push $(REPO):v$(COMMIT_HASH)
	sed -i "s/$(VERSION)/$(COMMIT_HASH)/g" chart/k8gb/Chart.yaml
endef

define debug
	$(call manifest)
	kubectl apply -f deploy/crds/test-namespace.yaml
	kubectl apply -f ./deploy/crds/k8gb.absa.oss_gslbs_crd.yaml
	kubectl apply -f ./deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml
	dlv $1
endef
