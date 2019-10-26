.PHONY: build
build:
	@[ -d .build ] || mkdir -p .build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o .build/app ./cmd/echo-service
	@file  .build/app
	@du -h .build/app

.PHONY: test
test:
	go test

.PHONY: deploy
deploy:
	@ansible-playbook deploy/deploy-service.yaml -i ./deploy/inventory
