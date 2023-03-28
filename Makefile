.PHONY: help build test generate clean

# Ref: https://gist.github.com/prwhite/8168133
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

GOMOCK_VERSION = "v1.5.0"

gomock: ## go generate mock file.
	go install "github.com/golang/mock/mockgen@$(GOMOCK_VERSION)"
	go list ./... |grep -v '/gomock' | xargs go generate -v

header: ## check and add license header.
	sh addlicense.sh

lint: ## run lint
ifeq (, $(shell which golangci-lint))
	# binary will be $(go env GOPATH)/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.51.2
else
	echo "Found golangci-lint"
endif
	golangci-lint run ./...


test: header lint ## Run test cases.
	go install "github.com/rakyll/gotest@v0.0.6"
	gotest -v -race -coverprofile=coverage.out -covermode=atomic ./...

e2e-test:
	go install "github.com/rakyll/gotest@v0.0.6"
	GIN_MODE=release
	LOG_LEVEL=fatal ## disable log for test
	gotest -v --tags=integration -race -coverprofile=coverage.out -covermode=atomic ./e2e/...

deps:  ## Update vendor.
	go mod verify
	go mod tidy -v

generate:  ## generate pb/tmpl file.
	# go get github.com/benbjohnson/tmpl
	# go install github.com/benbjohnson/tmpl
    # brew install flatbuffers
	sh ./proto/generate.sh

clean-mock: ## remove all mock files
	find ./ -name "*_mock.go" | xargs rm

clean-tmp: ## clean up tmp and test out files
	find . -type f -name '*.out' -exec rm -f {} +
	find . -type f -name '.DS_Store' -exec rm -f {} +
	find . -type f -name '*.test' -exec rm -f {} +
	find . -type f -name '*.prof' -exec rm -f {} +
	find . -type s -name 'localhost:*' -exec rm -f {} +
	find . -type s -name '127.0.0.1:*' -exec rm -f {} +

clean: clean-mock clean-tmp  ## Clean up useless files.
