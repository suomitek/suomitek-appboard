IMPORT_PATH:= github.com/kubeapps/kubeapps
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

all: kubeapps/dashboard kubeapps/apprepository-controller kubeapps/tiller-proxy kubeapps/kubeops kubeapps/assetsvc kubeapps/asset-syncer

# TODO(miguel) Create Makefiles per component
kubeapps/%:
	DOCKER_BUILDKIT=1 docker build -t kubeapps/$*$(IMG_MODIFIER):$(IMAGE_TAG) --build-arg "VERSION=${VERSION}" -f cmd/$*/Dockerfile .

kubeapps/dashboard:
	docker build -t kubeapps/dashboard$(IMG_MODIFIER):$(IMAGE_TAG) -f dashboard/Dockerfile dashboard/

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
