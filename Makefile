GO_CMD=go
GO_BUILD=$(GO_CMD) build
CMD_PATH=./cmd/myapp/main.go
DIST=dist
DIST_LINUX=$(DIST)/linux
BINARY_NAME=myapp
REPO=$(shell .github/scripts/reponame.sh)
DOCKER_BUILD_ARGS=($BUILD_ARGS)
GIT_TAG=$(TAG)

.PHONY: unit-test # Run unit-tests
unit-test:
	go test --short ./...

.PHONY: lint-test # Run lint-tests
lint-test:
	go mod tidy
	goimports -local github.com/macwis/go-boilerplate -w .
	go fmt ./...
	@golangci-lint -v run ./...
	@test -z "$$(golangci-lint run ./...)"

.PHONY: lint # Run lint-tests (alias)
lint: lint-test

.PHONY: arch-validate # Run clean architecture validator
arch-validate:
	@./.github/scripts/gocleanarch.sh
	cd internal/myapp
	go-cleanarch -ignore-tests -interfaces api -application usecase -domain domain

.PHONY: coverage # Generate coverage
coverage:
	@./.github.com/scripts/gocoverage.sh

.PHONY: build-linux # Build Linux binary
build-linux:
	mkdir -p $(DIST_LINUX)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o ./$(DIST_LINUX)/$(BINARY_NAME) -v $(CMD_PATH)

.PHONY: build # Build Linux binary (alias)
build: build-linux

.PHONY: build-image-application # Build Docker image
build-image-application:
	@docker build -f Dockerfile -t $(REGISTRY)/$(REPO):$(GIT_TAG) .

.PHONY: next-version # Identify next version (by tag)
next-version:
	@.github/scripts/nextversion.sh

.PHONY: last-tag # Identify current version (by tag)
last-tag:
	@.github/scripts/lasttag.sh

.PHONY: gen-tag # Push tag
gen-tag:
	@git tag $(GIT_TAG)
	@git push origin $(GIT_TAG)

.PHONY: application-name # Application name
application-name:
	@.github/scripts/reponame.sh

.PHONY: clean # Clean
clean:
	@docker images -a | grep '$(BINARY_NAME)' | awk '{print $3}' | xargs docker rmi --force || true
	rm ./$(DIST_LINUX)/$(BINARY_NAME) || true

.PHONY: wire-generate # Generate Wire bindings
wire-generate:
	cd internal/service/di ;\
	wire

.PHONY: help # Help - list of targets with descriptions
help:
	@echo ''
	@echo 'Usage:'
	@echo ' ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/ \1\t\2/' | expand -t20
