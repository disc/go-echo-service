# Requires goreleaser: https://goreleaser.com
.PHONY: build
build:
	goreleaser release --rm-dist --snapshot

.PHONY: test
test:
	go test

.PHONY: deploy
deploy: build
	ansible-playbook deploy/install-package.yaml -i ./deploy/inventory