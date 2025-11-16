.PHONY: dev-setup
dev-setup:
	go install github.com/go-delve/delve/cmd/dlv@latest							# for debugging
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest							# for generating type-safe database code
	go install github.com/swaggo/swag/cmd/swag@latest							# for generating swagger docs
	go install golang.org/x/tools/cmd/goimports@latest							# for formatting imports
	go install github.com/daixiang0/gci@latest									# for organizing imports
	go install mvdan.cc/gofumpt@latest											# for formatting code
	go install github.com/segmentio/golines@latest								# for formatting long lines
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest	# for linting and formatting

.PHONY: pre-commit
pre-commit: install generate format lint clean

# ============================================================================================================

.PHONY: install
install:
	go mod download

.PHONY: generate
generate:
# 	swag init -g http.go -o docs/swagger

.PHONY: lint
lint:
	golangci-lint run || true

.PHONY: format
format:
	golangci-lint run --fix || true
	gofumpt -l -w -extra .
	goimports -w .
	gci write \
		--custom-order -s standard -s default -s "prefix(sc-player-service)" -s blank \
		--no-lex-order --skip-generated --skip-vendor .
	golines -w -m 120 .

.PHONY: build
build: generate

.PHONY: dev
dev: build
	go run .

.PHONY: clean
clean:
	go mod tidy
	go mod verify
	go mod edit -fmt
	go clean
	rm -rf build/
