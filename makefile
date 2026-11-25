
.PHONY: help
help:
	@awk -F ':.*##' \
		' \
		BEGIN { printf "\nUsage:\n  make \033[1;93m<target>\033[0m\n" } \
		/^[a-zA-Z_0-9-]+:.*##/ { printf "  \033[1;93m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } \
		' \
		$(MAKEFILE_LIST)

##@ Development

.PHONY: test
test: manifests generate check ## Run all tests
	go test ./... -coverprofile cover.out

##@ Build

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: vendor
vendor: ## Run go mod vendor
	go mod vendor

.PHONY: build
build: ## Build the ns binary.
	go build -ldflags="-s -w" -o bin/ns ./cmd



