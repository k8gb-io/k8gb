REPO ?= ytsarev/ohmyglb
VERSION ?= $$(operator-sdk up local --operator-flags=-v)

.PHONY: up-local
up-local:
	kubectl create ns test-gslb
	kubectl apply -f ./deploy/crds/ohmyglb.absa.oss_gslbs_crd.yaml
	kubectl apply -f deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml
	operator-sdk up local --namespace=test-gslb

.PHONY: test
test:
	go test -v ./pkg/controller/gslb/

.PHONY: e2e-test
e2e-test:
	operator-sdk test local ./pkg/test/

.PHONY: dns-tools
dns-tools:
	kubectl -n ohmyglb get svc ohmyglb-coredns
	kubectl -n ohmyglb run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools

.PHONY: dns-smoke-test
dns-smoke-test:
	kubectl -n ohmyglb run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools --command -- /usr/bin/dig @gslb-coredns-coredns app3.cloud.absa.internal

.PHONY: deploy-local-cluster
deploy-local-cluster:
	kind create cluster --config=deploy/kind/cluster.yaml

.PHONY: destroy-local-cluster
destroy-local-cluster:
	kind delete cluster

.PHONY: deploy-gslb-operator
deploy-gslb-operator:
	kubectl apply -f deploy/namespace.yaml
	helm -n ohmyglb upgrade -i ohmyglb chart/ohmyglb

.PHONY: deploy-gslb-cr
deploy-gslb-cr:
	kubectl apply -f deploy/crds/test-namespace.yaml
	kubectl apply -f deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml

.PHONY: deploy-test-apps
deploy-test-apps:
	kubectl apply -f deploy/crds/test-namespace.yaml
	kubectl apply -f deploy/test-apps

.PHONY: clean-test-apps
clean-test-apps:
	kubectl delete -f deploy/test-apps

.PHONY: build
build:
	operator-sdk build $(REPO):$(VERSION)

.PHONY: push
push:
	docker push $(REPO):$(VERSION)
