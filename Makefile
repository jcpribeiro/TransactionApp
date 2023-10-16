.PHONY: help
help:  ## show this help
	@echo "usage: make [target]"
	@echo ""
	@egrep "^(.+)\:\ .*##\ (.+)" ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: run-services
run-services: ## run it will instance server
	docker-compose up -d

.PHONY: run
run: ## run it will instance server
	go run main.go

.PHONY: test
test: ## runing unit tests with covarage
	go test -v -failfast -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

.PHONY: mock
mock: ## mock is command to generate mock using mockgen
	go generate ./...

.PHONY: docs
docs: ## docs is a command to generate doc with swagger
	swag init --parseDependency --parseInternal --parseDepth 1

.PHONY: bump-deps
bump-deps: ## Update all dependencies
	go get -t -u ./...
	go mod tidy -compat=1.19

