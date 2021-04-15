HOSTNAME=softkraft.co
NAMESPACE=proj
NAME=nodeping
BINARY=terraform-provider-${NAME}
VERSION=0.1

ifndef OS_ARCH
	OS_ARCH=linux_amd64
endif


build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

run_tests: install
	go test -v -count=1 -timeout 30m ./nodeping
