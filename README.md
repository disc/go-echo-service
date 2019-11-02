# echo-service
Echo service is a Go lang application built based on go-kit approach and deployed using on goreleaser.

That's the test project for getting more knowledge about how to build and deploy microservices. 

### Requires
[goreleaser](https://goreleaser.com/introduction/)

### Engineering RPM build
```
goreleaser release --rm-dist --snapshot
```

### Ansible deploy
```
ansible-playbook deploy/install-package.yaml -i ./deploy/inventory
```