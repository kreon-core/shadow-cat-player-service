MODULE = sc-player-service

.PHONY: dev-setup
dev-setup:
	go install github.com/go-delve/delve/cmd/dlv@latest									# for debugging
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest									# for generating type-safe database code
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest	# for database migrations
# 	go install github.com/swaggo/swag/cmd/swag@latest									# for generating swagger docs
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest			# for linting and formatting
	go install golang.org/x/tools/cmd/goimports@latest									# for formatting imports
	go install github.com/daixiang0/gci@latest											# for organizing imports
	go install mvdan.cc/gofumpt@latest													# for formatting code
	go install github.com/segmentio/golines@latest										# for formatting long lines

.PHONY: pre-commit
pre-commit: install generate format lint tidy

# ============================================================================================================

.PHONY: install
install:
	go mod download

.PHONY: generate
generate:
	sqlc generate -f infrastructure/database/player/sqlc.yaml
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
		--custom-order -s standard -s default -s "prefix($(MODULE))" -s blank \
		--no-lex-order --skip-generated --skip-vendor .
	golines -w -m 120 .

.PHONY: build
build: generate
	go build -o build/app .

.PHONY: dev
dev: build
	./build/app -dev

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	go clean -cache -testcache -modcache
	rm -rf build

# ============================================================================================================

PLAYER_DB_DSN = postgres://postgres:password@localhost:5432/scs_player?sslmode=disable

.PHONY: migrate-player-create
migrate-player-create:
	migrate create -ext sql -dir infrastructure/database/player/migrations -seq $(name)

.PHONY: migrate-player-up
migrate-player-up:
	migrate -path infrastructure/database/player/migrations -database "$(PLAYER_DB_DSN)" up

.PHONY: migrate-player-down
migrate-player-down:
	migrate -path infrastructure/database/player/migrations -database "$(PLAYER_DB_DSN)" down
