SHELL := /bin/sh

GOBIN := $(shell go env GOPATH)/bin
GO_VERSION := $(shell go version | awk '{print $$3}')

PROTOLINT := $(GOBIN)/protolint
GOIMPORTS := $(GOBIN)/goimports
GOLANGCI_LINT := $(GOBIN)/golangci-lint
PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOBIN)/protoc-gen-go-grpc

PROTOC_GEN_GO_VERSION := v1.36.11
PROTOC_GEN_GO_GRPC_VERSION := v1.6.1
GOLANGCI_LINT_VERSION := v2.11.2

.PHONY: lint fix protobuf protolint gofmtcheck golint tools test build clean

lint: protolint gofmtcheck golint

fix:
	$(PROTOLINT) lint -fix .
	find . -type f -name '*.go' ! -name '*.pb.go' -print0 | xargs -0 $(GOIMPORTS) -w
	$(GOLANGCI_LINT) run --fix --timeout=2m -E gocritic -v ./...

protobuf:
	export PATH="$(GOBIN):$$PATH" && \
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(GOBIN)/protoc-gen-go \
		--plugin=protoc-gen-go-grpc=$(GOBIN)/protoc-gen-go-grpc \
		proto/proto/collector/collector.proto
	export PATH="$(GOBIN):$$PATH" && \
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(GOBIN)/protoc-gen-go \
		--plugin=protoc-gen-go-grpc=$(GOBIN)/protoc-gen-go-grpc \
		proto/proto/processor/processor.proto

protolint:
	$(PROTOLINT) .

gofmtcheck:
	find . -type f -name '*.go' ! -name '*.pb.go' -print0 | xargs -0 $(GOIMPORTS) -w
	git diff --exit-code

golint:
	$(GOLANGCI_LINT) run --timeout=2m -E gocritic -v ./...

tools:
	go install github.com/yoheimuta/protolint/cmd/protolint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

build:
	go build ./api/...
	go build ./collector/...
	go build ./processor/...

test:
	@echo "Running tests with $(GO_VERSION)"
	go test -race -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./api/... ./collector/... ./processor/... ./platform/...
	go tool cover -func=coverage.out | grep total
	@echo "✓ Tests completed. Coverage data in coverage.out"
	@echo "To view coverage details: go tool cover -html=coverage.out"

clean:
	find . -type f -name '*.pb.go' -delete
	find . -type f -name '*_grpc.pb.go' -delete