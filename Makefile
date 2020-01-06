.PHONY: e2e-tests
e2e-tests:
	kubectl apply -f deploy/namespace.yaml
	operator-sdk test local --namespace ohmyglb ./pkg/test/
