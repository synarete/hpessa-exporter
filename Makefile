#
# Copyright 2022 Red Hat OpenShift Data Foundation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Control variables.
# Developer may define custom variables via optional 'hack/make-devel-vars.mk'
-include devel-vars.mk
include bundle-vars.mk

# Helper functions: report executed target
define report
	$(info $(word 1, $(MAKEFILE_LIST)): $(notdir $@))
endef

# Be verbode with 'make V=1' or quite by default
Q = @
ifeq ("$(origin V)", "command line")
ifeq ($(V), 1)
Q =
endif
endif

# Follow some kubebuilder conventions
# See: https://book.kubebuilder.io/cronjob-tutorial/basic-project.html
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Directory layout variables
PROJECT_DIR := $(CURDIR)
HACK_DIR := $(PROJECT_DIR)/hack
BUILD_DIR := $(PROJECT_DIR)/build/
BUILD_TMP_DIR := $(BUILD_DIR)/tmp/
TOOLS_BIN_DIR := $(BUILD_DIR)/tools/bin/
OUTPUT_DIR := $(BUILD_DIR)/output/
OUTPUT_BIN_DIR := $(OUTPUT_DIR)/bin/
OUTPUT_CFG_DIR := $(OUTPUT_DIR)/config/

# Common build variables
GITREV := $(shell $(HACK_DIR)/gitrev.sh)
RM := rm -rf
MKDIR_P := mkdir -p
GO ?= $(shell command -v go)
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
GOBIN ?= $(TOOLS_BIN_DIR)
GO111MODULE=on

# Go sources variables
GOPREF := github.com/red-hat-storage/hpessa-exporter
GODDIR := $(GOPREF)/internal/devmon
GOVARS := -X '$(GODDIR).version=$(VERSION)' -X '$(GODDIR).gitrev=$(GITREV)'
GOSRCS := $(shell find cmd internal -name \*.go)

# Extra build tools
KUSTOMIZE ?= $(GOBIN)/kustomize
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
YQ ?= $(GOBIN)/yq

# Container image & deployment tools
DOCKERCMD ?= podman
KUBECTL ?= kubectl


.DEFAULT_GOAL := all
.PHONY: all
all: version build

##@ General

.PHONY: help
help: ## Display this help.
	$(Q)$(HACK_DIR)/make-help.awk $(MAKEFILE_LIST)

.PHONY: version
version:  ## Show version info
	$(Q)$(call report, $@)
	$(Q)echo $(VERSION)-$(GITREV)
	$(Q)$(GO) version

##@ Development
.PHONY: pre-build
pre-build: pre-build-tools pre-build-layout

.PHONY: pre-build-tools
pre-build-tools:
	$(Q)$(GO) version > /dev/null
	$(Q)$(DOCKERCMD) version > /dev/null

.PHONY: pre-build-layout
pre-build-layout:
	$(Q)$(MKDIR_P) $(BUILD_TMP_DIR)
	$(Q)$(MKDIR_P) $(TOOLS_BIN_DIR)
	$(Q)$(MKDIR_P) $(OUTPUT_BIN_DIR)
	$(Q)$(MKDIR_P) $(OUTPUT_CFG_DIR)


.PHONY: build
build: build-tools go-deps-update go-fmt go-vet go-build ## Build program

.PHONY: go-deps-update
go-deps-update:
	$(Q)$(call report, $@)
	$(Q)cd $(PROJECT_DIR) && $(GO) mod tidy
	$(Q)cd $(PROJECT_DIR) && $(GO) mod vendor

.PHONY: yaml-fmt
yaml-fmt: install-yq ## Generate consistent yaml files
	$(Q)$(call report, $@)
	$(Q)YQ=$(YQ) find $(PROJECT_DIR)/config -type f -name '*.yaml' \
		-exec $(HACK_DIR)/yq-fixup-files.sh {} \;

.PHONY: go-fmt
go-fmt:
	$(Q)$(call report, $@)
	$(Q)cd $(PROJECT_DIR) && $(GO) fmt ./...

.PHONY: go-vet
go-vet:
	$(Q)$(call report, $@)
	$(Q)cd $(PROJECT_DIR) && $(GO) vet ./...

PHONY: go-mod
go-mod:
	$(Q)$(call report, $@)
	$(Q)cd $(PROJECT_DIR) && $(GO) mod tidy
	$(Q)cd $(PROJECT_DIR) && $(GO) mod vendor

.PHONY: go-build
go-build: $(OUTPUT_BIN_DIR)/hpessa-exporter

.PHONY: go-srcs
go-srcs: $(GOSRCS)

$(OUTPUT_BIN_DIR)/hpessa-exporter: go-srcs
	$(Q)$(call report, $@)
	$(Q)cd $(PROJECT_DIR) && \
		CGO_ENABLED=0 GOOS=$(GOOS) GOBIN=$(GOBIN) \
		$(GO) build -a -o $@ \
		-ldflags="-s -w $(GOVARS)" $(PROJECT_DIR)/cmd/main.go

.PHONY: lint
lint: build-tools go-lint ## Lint Go sources

.PHONY: golangci-lint
go-lint:
	$(Q)$(call report, $@)
	$(Q)GOLANGCI_LINT_CACHE=$(BUILD_TMP_DIR)/golangci-cache \
		$(GOLANGCI_LINT) -c $(PROJECT_DIR)/.golangci.yaml run ./...

.PHONY: test
test: go-test ## Run test suite

go-test:
	$(Q)$(call report, $@)
	$(Q)cd $(PROJECT_DIR) && \
		$(GO) test -v -failfast -p 1 -cover $(shell $(GO) list ./...)

.PHONY:
clean: ## Clean build outputs
	$(Q)$(call report, $@)
	$(Q)$(RM) $(BUILD_TMP_DIR)/
	$(Q)$(RM) $(OUTPUT_DIR)

##@ Build tools

.PHONY: build-tools
build-tools: pre-build-layout install-yq install-kustomize install-golangci-lint

.PHONY: install-yq
install-yq: $(YQ) ## Download and install yq

$(YQ):
	$(Q)$(call report, $@)
	$(Q)$(MKDIR_P) $(BUILD_TMP_DIR)/yq
	$(Q)GOBIN=$(GOBIN) $(HACK_DIR)/install-yq.sh $@ $(BUILD_TMP_DIR)/yq
	$(Q)$(RM) $(BUILD_TMP_DIR)/yq

.PHONY: install-kustomize
install-kustomize: $(KUSTOMIZE) ## Download and install kustomize

$(KUSTOMIZE):
	$(Q)$(call report, $@)
	$(Q)$(MKDIR_P) $(BUILD_TMP_DIR)/kustomize
	$(Q)GOBIN=$(GOBIN) $(HACK_DIR)/install-kustomize.sh $@ \
		$(BUILD_TMP_DIR)/kustomize
	$(Q)$(RM) $(BUILD_TMP_DIR)/kustomize

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT) ## Download and install golangci-lint

$(GOLANGCI_LINT):
	$(Q)$(call report, $@)
	$(Q)$(HACK_DIR)/install-golangci-lint.sh $@

##@ Container image

.PHONY: image
image: build-tools image-build image-push ## Build and push container image

.PHONY: image-build
image-build: Dockerfile ## Build and tag container image
	$(Q)$(call report, $@)
	$(Q)$(DOCKERCMD) build -t $(IMG) -f $< \
		--build-arg ARCH=$(GOARCH) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GITREV=$(GITREV) \
		--build-arg GOVARS="$(GOVARS)"

.PHONY: image-push
image-push: ## Push container image to remote repository
	$(Q)$(call report, $@)
	$(Q)$(DOCKERCMD) push $(IMG)

##@ Deployment

.PHONY: deploy
deploy: build-tools ## Deploy controller to cluster
	$(Q)$(call report, $@)
	$(Q)cp -r config/* $(OUTPUT_CFG_DIR)/
	$(Q)cd $(OUTPUT_CFG_DIR) && $(KUSTOMIZE) edit set image hpessa-exporter=$(IMG)
	$(Q)$(KUSTOMIZE) build $(OUTPUT_CFG_DIR) | $(KUBECTL) apply -f -

.PHONY: undeploy
undeploy: build-tools ## Undeploy controller from the cluster
	$(Q)$(KUSTOMIZE) build $(OUTPUT_CFG_DIR) | $(KUBECTL) delete -f -
