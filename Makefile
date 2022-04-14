GOFILES = $(shell go list -mod vendor ./... | grep -v vendor)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: analyze
analyze: lint vet test ## Run lint, vet, and test

.PHONY: lint
lint: ## Lint the code
    # TODO: revive seems to give an error exit code even though it doesn't report any problems. How to handle?
	@revive -config .revive.toml ./... | grep -v vendor

.PHONY: test
test: unit-test ## Run the test suite(s).

.PHONY: unit-test
unit-test: ## Run the unit tests.
	@gotestsum --format=standard-verbose -- -run '(?i)unit' ./...

.PHONY: vet
vet: ## Verify `go vet` passes.
	@go vet -mod vendor $(GOFILES)

.PHONY: setup-osx-env
setup-osx-env: brew-deps setup-asdf install-goimports install-gotestsum asdf-reshim setup-pre-commit ## Setup the development environment for OS X

.PHONY: brew-deps
brew-deps: ## Install tools via Homebrew
	@echo ">>>> Installing tools"
	@brew bundle install --no-lock --file Brewfile

.PHONY: setup-asdf
setup-asdf: ## Install Go via asdf
	@echo ">>>> Installing asdf plugins"
	@asdf plugin add golang
	@echo ">>>> Installing Go via asdf"
	@asdf install

.PHONY: install-goimports
install-goimports: ## Install goimports
	@echo ">>>> Installing goimports"
	@go get golang.org/x/tools/cmd/goimports

.PHONY: install-gotestsum
install-gotestsum: ## Install gotestsum
	@echo ">>>> Installing gotestsum"
	@go get gotest.tools/gotestsum

.PHONY: asdf-reshim
asdf-reshim: ## Reshim asdf
	@asdf reshim golang

.PHONY: setup-pre-commit
setup-pre-commit: ## Setup pre-commit hooks
	@echo ">>>> Setting up pre-commit"
	@git config --global init.templateDir ~/.git-template
	@pre-commit init-templatedir -t pre-commit -t prepare-commit-msg ~/.git-template
	@pre-commit install --install-hooks --allow-missing-config -t pre-commit -t prepare-commit-msg -f