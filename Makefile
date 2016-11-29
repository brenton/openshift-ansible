.PHONY: test
test:
	nosetests

.PHONY: test-integration
test-integration:
	# TODO: remove `sudo` when the CI environment can run docker CLI commads without it.
	sudo go test -v ./test/integration/...
