REPO ?= ytsarev/ohmyglb
VERSION ?= $$(operator-sdk up local --operator-flags=-v)
VALUES_YAML ?= chart/ohmyglb/values.yaml
ETCD_DEBUG_IMAGE ?= quay.io/coreos/etcd:v3.2.25
GSLB_DOMAIN ?= cloud.example.com

.PHONY: up-local
up-local: create-test-ns
	kubectl apply -f ./deploy/crds/ohmyglb.absa.oss_gslbs_crd.yaml
	kubectl apply -f deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml
	operator-sdk up local --namespace=test-gslb

.PHONY: lint
lint:
	staticcheck ./pkg/...
	errcheck ./pkg/...

.PHONY: test
test:
	go test -v ./pkg/controller/gslb/

.PHONY: e2e-test
e2e-test: deploy-gslb-operator
	operator-sdk test local ./pkg/test --no-setup --namespace ohmyglb

.PHONY: dns-tools
dns-tools:
	kubectl -n ohmyglb get svc ohmyglb-coredns
	kubectl -n ohmyglb run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools

.PHONY: dns-smoke-test
dns-smoke-test:
	kubectl -n ohmyglb run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools --command -- /usr/bin/dig @ohmyglb-coredns app3.cloud.example.com

.PHONY: deploy-local-cluster
deploy-local-cluster:
	kind create cluster --config=deploy/kind/cluster.yaml

.PHONY: destroy-local-cluster
destroy-local-cluster:
	kind delete cluster

.PHONY: create-ohmyglb-ns
create-ohmyglb-ns:
	kubectl apply -f deploy/namespace.yaml

.PHONY: create-test-ns
create-test-ns:
	kubectl apply -f deploy/crds/test-namespace.yaml

.PHONY: deploy-local-ingress
deploy-local-ingress: create-ohmyglb-ns
	helm -n ohmyglb upgrade -i nginx-ingress stable/nginx-ingress -f deploy/ingress/nginx-ingress-values.yaml

.PHONY: deploy-gslb-operator
deploy-gslb-operator: create-ohmyglb-ns
	cd chart/ohmyglb && helm dependency update
	helm -n ohmyglb upgrade -i ohmyglb chart/ohmyglb -f $(VALUES_YAML)

# workaround until https://github.com/crossplaneio/crossplane/issues/1170 solved
.PHONY: deploy-gslb-operator-14
deploy-gslb-operator-14: create-ohmyglb-ns
	cd chart/ohmyglb && helm dependency update
	helm -n ohmyglb template ohmyglb chart/ohmyglb -f $(VALUES_YAML) | kubectl -n ohmyglb --validate=false apply -f -

.PHONY: deploy-gslb-cr
deploy-gslb-cr: create-test-ns
	sed -i 's/cloud\.example\.com/$(GSLB_DOMAIN)/g' deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml
	kubectl apply -f deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml
	git checkout -- deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml

.PHONY: deploy-test-apps
deploy-test-apps: create-test-ns
	kubectl apply -f deploy/test-apps
	helm repo add podinfo https://stefanprodan.github.io/podinfo
	helm upgrade --install --wait frontend --namespace test-gslb -f deploy/test-apps/podinfo/podinfo-values.yaml --set backend=http://backend-podinfo:9898/echo podinfo/podinfo
	helm upgrade --install --wait backend --namespace test-gslb -f deploy/test-apps/podinfo/podinfo-values.yaml podinfo/podinfo

.PHONY: clean-test-apps
clean-test-apps:
	kubectl delete -f deploy/test-apps
	helm -n test-gslb uninstall backend
	helm -n test-gslb uninstall frontend

.PHONY: build
build:
	operator-sdk build $(REPO):$(VERSION)

.PHONY: push
push:
	docker push $(REPO):$(VERSION)

.PHONY: debug-test-etcd
debug-test-etcd:
	kubectl run --rm -i --tty --env="ETCDCTL_API=3" --env="ETCDCTL_ENDPOINTS=http://etcd-cluster-client:2379" --namespace ohmyglb etcd-test --image "$(ETCD_DEBUG_IMAGE)" --restart=Never -- /bin/sh

.PHONY: infoblox-secret
infoblox-secret:
	kubectl -n ohmyglb create secret generic external-dns \
	    --from-literal=EXTERNAL_DNS_INFOBLOX_WAPI_USERNAME=$${WAPI_USERNAME} \
	    --from-literal=EXTERNAL_DNS_INFOBLOX_WAPI_PASSWORD=$${WAPI_PASSWORD}
