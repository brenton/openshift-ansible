.PHONY: test
test:
	nosetests

.PHONY: test-integration
test-integration:
	sudo go test -v ./test/integration/...
