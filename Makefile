GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
EXPORT_RESULT ?= false
GO_TOOLS = github.com/vektra/mockery/v2/.../
ENV_BUILD_LOCAL_DOCKER=APP_IMAGE_NAME=git-watchdog:local
ENV_INTEGRATION_TEST=\
  DB_ADDRESS=0.0.0.0:3307 \
  DB_PASSWORD=test \
  DB_USER=test \
  DB_NAME=code-challenge

BUF_VERSION:=0.55.0
.PHONY: install install-go-tools install-osx lint all test vendor coverage

all: help

help: ## Show this help.
ifeq ($(EXPORT_RESULT),true)
	@echo "OK"
else
	@echo "Fail"
endif
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)

install-go-tools: ## Install go dependencies
	go get $(GO_TOOLS)

install:
	curl -sSL \
    	"https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" \
    	-o "$(shell go env GOPATH)/bin/buf" && \
  	chmod +x "$(shell go env GOPATH)/bin/buf"

generate:
	buf generate

lint: lint-go

docker.integration.start:
	docker-compose -f docker-compose-it.yml up -d --remove-orphans mydb;

docker.integration.stop:
	docker-compose -f docker-compose-it.yml down -v;

docker.local.start:
	$(ENV_BUILD_LOCAL_DOCKER) docker-compose -f docker-compose.yml up -d --build --remove-orphans;

docker.local.stop:
	$(ENV_BUILD_LOCAL_DOCKER) docker-compose -f docker-compose.yml down -v;

build:
	@./scripts/app-build.sh

build.docker.image:
	@docker build -t git-watchdog:local .

migrate: 
	@./scripts/migrate.sh

build.migrate:
	@./scripts/migrate-build.sh

build.consumer:
	@./scripts/scanworker-build.sh

build.scanworker:
	@./scripts/scanworker-build.sh

lint-go: ## Linting go files
	@golangci-lint run --allow-parallel-runners

run:
	@go run cmd/app/main.go

clean: ## Remove
	rm -f ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov

test.unit: ## Run the tests of the project
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

test.integration:
	$(ENV_INTEGRATION_TEST) $(GOTEST) -tags=integration ./internal/integration -v -count=1

coverage: ## Generate coverage report
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/AlekSi/gocov-xml
	GO111MODULE=off go get -u github.com/axw/gocov/gocov
	gocov convert profile.cov | gocov-xml > coverage.xml
endif
