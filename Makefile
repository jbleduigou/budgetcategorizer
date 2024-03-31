.DEFAULT_GOAL := help

check: test ## Runs all tests

test: ## Run the unit tests
	go test -cover -v $(shell go list ./... | grep -v /vendor/)

clean: ## Clean up build artifacts
	go clean
	rm -f bootstrap bootstrap.zip

build: clean test ## Build the executable
	GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -tags lambda.norpc -o ./bootstrap ./cmd/budgetcategorizer

zip: build ## Zip the executable so that it can be uploaded to AWS Lambda
	zip bootstrap.zip bootstrap

build-BudgetCategorizerFunction:
	go get -v -t -d ./...
	GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -tags lambda.norpc -o ./cmd/budgetcategorizer/bootstrap ./cmd/budgetcategorizer
	cp ./cmd/budgetcategorizer/bootstrap $(ARTIFACTS_DIR)/bootstrap

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.SILENT: zip build test lint vet clean help

