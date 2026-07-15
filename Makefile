VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)
BUILD_FLAGS := -trimpath -buildvcs=false -ldflags "$(LDFLAGS)"

.PHONY: build test vet verify release

build:
	go build $(BUILD_FLAGS) -o cairn ./cmd/cairn

test:
	go test ./...

vet:
	go vet ./...

verify: test vet build

release:
	@test "$(VERSION)" != "dev" || (echo "VERSION e' obbligatoria, es. make release VERSION=v0.1.0" >&2; exit 1)
	mkdir -p dist
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/cairn_$(VERSION)_darwin_amd64 ./cmd/cairn
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/cairn_$(VERSION)_darwin_arm64 ./cmd/cairn
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/cairn_$(VERSION)_linux_amd64 ./cmd/cairn
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/cairn_$(VERSION)_linux_arm64 ./cmd/cairn
