NAMESPACE=terraform
HOSTNAME=softkraft.co
NAME=nodeping
BINARY=out/terraform-provider-${NAME}
VERSION=0.0.1

OS_ARCH=$(shell terraform -v | sed -n  '/on /p' | sed -e 's/\on //g')
TF_SUPPORTED_VERSION=1.0.0


TF_VER_NUM=$(shell terraform -v | sed -n  '/Terraform v/p' | sed -e 's/\Terraform v//g')
TF_VER_MAJOR=$(shell echo $(TF_VER_NUM) | cut -f1 -d.)
TF_VER_MINOR=$(shell echo $(TF_VER_NUM) | cut -f2 -d.)
TF_VER_PATCH=$(shell echo $(TF_VER_NUM) | cut -f3 -d.)

check_terraform:
ifeq (, $(shell which terraform))
	$(error "Cannot find terraform, please visit https://www.terraform.io/")
endif
	$(eval TF_SUP_MAJOR := $(shell echo $(TF_SUPPORTED_VERSION) | cut -f1 -d.))
	$(eval TF_SUP_MINOR := $(shell echo $(TF_SUPPORTED_VERSION) | cut -f2 -d.))
	$(eval TF_SUP_PATCH := $(shell echo $(TF_SUPPORTED_VERSION) | cut -f3 -d.))

	@if [ $(TF_VER_MAJOR) -ge $(TF_SUP_MAJOR) ] && [ $(TF_VER_MINOR) -ge $(TF_SUP_MINOR) ] && [ $(TF_VER_PATCH) -ge $(TF_SUP_PATCH) ]; then \
        echo "Required terraform version has found"; \
    else \
		echo "Wrong terraform version. Minimal supported version: 1.0.0"; exit 1; \
    fi

all: vendoring fmt build

fmt:
	go fmt main.go
	go fmt nodeping/*.go
	go fmt nodeping_api_client/*.go

vendoring:
	go mod vendor

build:
	go build -o ${BINARY}

install: vendoring fmt build check_terraform
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

run_tests: install
	go test -v -count=1 -timeout 30m ./nodeping
