.PHONY: test
test:
	go test -v ./pkg/controller/gslb/

.PHONY: e2e-test
e2e-test:
	operator-sdk test local ./pkg/test/
