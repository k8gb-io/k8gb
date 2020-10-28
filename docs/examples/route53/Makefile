CLUSTER_NAME ?= k8gb-cluster-eu-west-1

.PHONY: login
login:
	aws eks --region eu-west-1 update-kubeconfig --name $(CLUSTER_NAME)

.PHONY: k8gb-reference-setup
k8gb-reference-setup:
	terraform apply
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.40.2/deploy/static/provider/aws/deploy.yaml
