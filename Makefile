.PHONY: swag fmt-go gen-ent install-atlas gen-migration apply-migration-test hash-migration-test apply-migration-prod

USER := flashdeal
HOST_TEST := localhost
DB_URL_TEST := "mysql://flashdeal:FwT93%23bzC2%26jLpNq@localhost:3306/flashdeal_db?charset=utf8mb4&parseTime=True&loc=UTC"
DB_URL_PROD := mysql://admin:pass@host:3306/flashdeal_db
COMMIT := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
APP_VERSION ?= v0.1.0
APP_SUFFIX ?= test
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

swag:
	@echo "> Generating Swagger docs"
	swag fmt && swag init -g=./cmd/apiserver/main.go -o=./api --parseDependency --parseInternal

fmt-go:
	@echo "> Formatting Go code"
	go fmt ./...
	goimports -v --format-only -w -local=github.com/algorithm9/flash-deal ./cmd ./internal ./pkg

gen-ent:
	@echo "> Generating Ent code"
	go generate ./internal/pkg/ent

install-atlas:
	@echo "> Installing Atlas (if needed)"
	go install ariga.io/atlas/cmd/atlas@latest

gen-migration:
	@echo "> Generating Atlas migration file"
	atlas migrate diff $(diffname) \
      --dir "file://internal/shared/ent/migrations" \
      --to "ent://internal/shared/ent/schema" \
      --dev-url "docker://mysql/8/ent" \
      --format '{{ sql . "    " }}'

apply-migration-test:
	@echo "> Applying migration to test DB"
	atlas migrate apply \
	   --dir "file://internal/shared/ent/migrations" \
	   --url "$(DB_URL_TEST)"

status-migration-test:
	atlas migrate status \
	   --dir "file://internal/shared/ent/migrations" \
	   --url "$(DB_URL_TEST)"

hash-migration-test:
	atlas migrate hash \
		--dir "file://internal/shared/ent/migrations"

apply-migration-prod:
	atlas migrate apply \
	   --dir "file://internal/shared/ent/migrations" \
	   --url "$(DB_URL_PROD)"

build-apiserver-binary: ##@ Build flashdeal apiserver
	@printf "==> Begin building flashdeal-apiserver, os=%s, arch=%s, time='%s', version=%s\n" "$(GOOS)" "$(GOARCH)" "$(BUILD_TIME)" "$(APP_VERSION)"
	@start=$$(date +%s); \
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X 'github.com/algorithm9/flash-deal/meta.BuildTime=$(BUILD_TIME)' -X 'github.com/algorithm9/flash-deal/meta.Version=$(APP_VERSION)' -X 'github.com/algorithm9/flash-deal/meta.Commit=$(COMMIT)'" -o bin/flashdeal-apiserver$(APP_SUFFIX) ./cmd/apiserver; \
	end=$$(date +%s); \
	printf "🍺  Build completed in %s seconds:\n" $$((end - start))
	@ls -lFh bin/flashdeal-apiserver$(APP_SUFFIX)
