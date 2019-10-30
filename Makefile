.PHONY: build
build:
	@[ -d .build ] || mkdir -p .build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o .build/echo-service ./cmd/echoservice/main.go
	@file  .build/echo-service
	@du -h .build/echo-service

.PHONY: test
test:
	go test

.PHONY: deploy
deploy:
	@ansible-playbook deploy/deploy-service.yaml -i ./deploy/inventory

# Requires install of https://github.com/goreleaser/nfpm
build-rpm:
	nfpm pkg -f deploy/service.yaml --target deploy/service.rpm

install-rpm: build-rpm
	ansible-playbook deploy/install-package.yaml -i ./deploy/inventory