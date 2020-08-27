IMPORT_PATH:= github.com/suomitek/suomitek-appboard
GO = /usr/bin/env go
GOFMT = /usr/bin/env gofmt
IMAGE_TAG ?= dev-$(shell date +%FT%H-%M-%S-%Z)
VERSION ?= $$(git rev-parse HEAD)

default: all

include ./script/cluster-kind.mk
include ./script/cluster-openshift.mk
include ./script/deploy-dev.mk

IMG_MODIFIER ?=

GO_PACKAGES = ./...
# GO_FILES := $(shell find $(shell $(GO) list -f '{{.Dir}}' $(GO_PACKAGES)) -name \*.go)

all: suomitek-appboard/dashboard suomitek-appboard/apprepository-controller suomitek-appboard/tiller-proxy suomitek-appboard/kubeops suomitek-appboard/assetsvc suomitek-appboard/asset-syncer

# TODO(miguel) Create Makefiles per component
suomitek-appboard/%:
	DOCKER_BUILDKIT=1 docker build -t suomitek-appboard/$*$(IMG_MODIFIER):$(IMAGE_TAG) --build-arg "VERSION=${VERSION}" -f cmd/$*/Dockerfile .

suomitek-appboard/dashboard:
	docker build -t suomitek-appboard/dashboard$(IMG_MODIFIER):$(IMAGE_TAG) -f dashboard/Dockerfile dashboard/

test:
	$(GO) test $(GO_PACKAGES)

test-db:
	# It's not supported to run tests that involve a database in parallel since they are currently
	# using the same PG schema. We need to run them sequentially 
	cd cmd/asset-syncer; ENABLE_PG_INTEGRATION_TESTS=1 go test -count=1 ./...
	cd cmd/assetsvc; ENABLE_PG_INTEGRATION_TESTS=1 go test -count=1 ./...

test-all: test-apprepository-controller test-dashboard

test-dashboard:
	yarn --cwd dashboard/ install --frozen-lockfile
	yarn --cwd=dashboard run lint
	CI=true yarn --cwd dashboard/ run test

test-%:
	$(GO) test -v $(IMPORT_PATH)/cmd/$*

fmt:
	$(GOFMT) -s -w $(GO_FILES)

vet:
	$(GO) vet $(GO_PACKAGES)

.PHONY: default all test-all test test-dashboard fmt vet
