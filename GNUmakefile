NAME=outscale
BINARY=packer-plugin-${NAME}
PLUGIN_FQN="$(shell grep -E '^module' <go.mod | sed -E 's/module *//')"

COUNT?=1
TEST?=$(shell go list ./...)
CHECK_DIR?=./...
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

.PHONY: fmt

.PHONY: dev

.PHONY: vet

fmt:
	gofmt -l -s -w .

vet: fmt
	go vet $(CHECK_DIR)

build: fmt vet
	go build -o ${BINARY}

dev:
	@go build -ldflags="-X '${PLUGIN_FQN}/version.VersionPrerelease=dev'" -o '${BINARY}'
	packer plugins install --path ${BINARY} "$(shell echo "${PLUGIN_FQN}" | sed 's/packer-plugin-//')"

test:
	go test -race -count $(COUNT) $(TEST) -timeout=3m

plugin-check: build
	go run -modfile=./go.mod github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc plugin-check ${BINARY}

testacc: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m

generate:
	go generate ./...
	if [ -d ".docs" ]; then rm -r ".docs"; fi
	go run -modfile=./go.mod github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc renderdocs -src "docs" -partials docs-partials/ -dst ".docs/"
	./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs" "outscale"
	rm -r ".docs"
	# checkout the .docs folder for a preview of the docs
