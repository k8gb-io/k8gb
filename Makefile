.PHONY: e2e-tests
e2e-tests:
	operator-sdk test local ./pkg/test/
