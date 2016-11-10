.PHONY: test
test:
	nosetests

.PHONY: test-integration
test-integration:
	go test -v ./test/integration/...
