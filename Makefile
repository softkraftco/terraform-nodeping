HOSTNAME=softkraft.co
NAMESPACE=terraform
NAME=nodeping
BINARY=out/terraform-provider-${NAME}
VERSION=0.0.1

ifndef OS_ARCH
	OS_ARCH=linux_amd64
endif

all: vendor build

fmt:
	go fmt main.go
	go fmt nodeping/*.go

vendoring:
	go mod vendor

build:
	go build -o ${BINARY}

install: fmt vendoring build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

run_tests: install
	go test -v -count=1 -timeout 30m ./nodeping
