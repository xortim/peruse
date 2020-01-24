GO111MODULE  = on
GOPATH      ?= ${HOME}/go
DIST        := ./dist
EXECUTABLE  := peruse
GITVERSION  := $(shell git describe --dirty --always --tags --long)
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
GOBUILD     := go build -trimpath -ldflags "-X ${PACKAGENAME}/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/conf.GitVersion=${GITVERSION}"
LINUXOUT    := ${EXECUTABLE}-linux

KUBECONFIG  ?= $(shell k3d get-kubeconfig --name='k3s-default')

default: ${EXECUTABLE}

.PHONY: all default clean build test deps create-cluster

${EXECUTABLE}: dist deps
	$(GOBUILD) -o $(DIST)/$(EXECUTABLE)

build-linux:
	GOOS=linux GOOARCH=amd64 $(GOBUILD) -o $(DIST)/$(LINUXOUT)

test:
	go test -cover ./...

deps:
	go get -d -v # Adding -u here will break CI

clean:
	go clean -modcache

dist:
	@mkdir -p $(DIST)

create-cluster:
	@k3d create --publish 10080:80 --publish 10443:443 --workers 2
	@export KUBECONFIG="$(k3d get-kubeconfig --name='k3s-default')"

