CLIENT_SOURCE=./cmd/ligoloc
SERVER_SOURCE=./cmd/ligolos
TLS_CERT ?= 'certs/cert.pem'
LDFLAGSOLD="-s -w -X main.tlsFingerprint=$$(openssl x509 -fingerprint -sha256 -noout -in $(TLS_CERT) | cut -d '=' -f2)"
LDFLAGS="-s -w"
GCFLAGS="all=-trimpath=$GOPATH"

CLIENT_BINARY=ligoloc
SERVER_BINARY=ligolos
TAGS=release

OSARCH = "linux/amd64 linux/386 windows/amd64 windows/386 darwin/amd64 darwin/386"

TLS_HOST ?= 'ligolo.lan'

.DEFAULT: help

help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dep: ## Install dependencies
	go get -d -v ./...
	go get -u github.com/mitchellh/gox

certs: ## Build SSL certificates
	mkdir certs
	cd certs && go run `go env GOROOT`/src/crypto/tls/generate_cert.go -ecdsa-curve P256 -ed25519 -host $(TLS_HOST)

build: ## Build for the current architecture.
	go build -ldflags $(LDFLAGS) -gcflags $(GCFLAGS) -tags $(TAGS) -o bin/$(CLIENT_BINARY) $(CLIENT_SOURCE) && \
	go build -ldflags $(LDFLAGS) -gcflags $(GCFLAGS) -tags $(TAGS) -o bin/$(SERVER_BINARY) $(SERVER_SOURCE)

build-all: ## Build for every architectures.
	gox -osarch=$(OSARCH) -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -tags $(TAGS) -output "bin/$(SERVER_BINARY)_{{.OS}}_{{.Arch}}" $(SERVER_SOURCE)
	gox -osarch=$(OSARCH) -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -tags $(TAGS) -output "bin/$(CLIENT_BINARY)_{{.OS}}_{{.Arch}}" $(CLIENT_SOURCE)

clean:
	rm -rf certs
	rm bin/$(SERVER_BINARY)_*
	rm bin/$(CLIENT_BINARY)_*
