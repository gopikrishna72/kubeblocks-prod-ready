GO ?= go
GOPATH ?= $(shell $(GO) env GOPATH)

GINKGO_VERSION = 2.9.1
GINKGO = $(GOPATH)/bin/ginkgo
HELM_VERSION ?= v3.9.0
INIT_ENV ?= false
CONFIG_TYPE ?= s3
# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html



.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY:ginkgo
ginkgo: # Make sure ginkgo is in $GOPATH/bin
	go install github.com/onsi/ginkgo/v2/ginkgo@v$(GINKGO_VERSION)

ifeq ($(origin VERSION), command line)
    VERSION ?= $(VERSION)
endif
ifeq ($(origin PROVIDER), command line)
    PROVIDER ?= $(PROVIDER)
else
    PROVIDER ?= ""
endif
ifeq ($(origin REGION), command line)
    REGION ?= $(REGION)
else
    REGION ?= ""
endif
ifeq ($(origin SKIP_CASE), command line)
    SKIP_CASE ?= $(SKIP_CASE)
else
    SKIP_CASE ?= ""
endif
ifeq ($(origin INIT_ENV), command line)
    INIT_ENV ?= $(INIT_ENV)
else
    INIT_ENV ?= false
endif
ifeq ($(origin CONFIG_TYPE), command line)
    CONFIG_TYPE ?= $(CONFIG_TYPE)
else
    CONFIG_TYPE ?= s3
endif

.PHONY: run
run: ginkgo ## Run end-to-end tests.
	#ACK_GINKGO_DEPRECATIONS=$(GINKGO_VERSION) $(GINKGO) run .
	$(GINKGO) test -process -ginkgo.v . -- -VERSION=$(VERSION) -PROVIDER=$(PROVIDER) -REGION=$(REGION) -SECRET_ID=$(SECRET_ID) -SECRET_KEY=$(SECRET_KEY) -INIT_ENV=$(INIT_ENV) -TEST_TYPE=$(TEST_TYPE) -SKIP_CASE=$(SKIP_CASE) -CONFIG_TYPE=$(CONFIG_TYPE) --ginkgo.json-report=report.json

build: ginkgo ## Run ginkgo build e2e test suite binary.
	$(GINKGO) build .

.PHONY: helmtool
helmtool: ## Download helm locally if necessary.
ifeq (, $(shell which helm))
	@{ \
	set -e ;\
	go install github.com/helm/helm@$(HELM_VERSION);\
	}
HELM=$(GOBIN)/helm
else
HELM=$(shell which helm)
endif