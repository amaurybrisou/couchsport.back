ifndef GOROOT 
$(error "GOROOT is not set")
endif

ifndef GOBIN
$(error "GOBIN is not set")
endif

ENV=local

RELEASE_PATH=$(CURDIR)/release

GOCMD=$(GOROOT)/bin/go
DOCKERCMD=$(shell which docker)

.DEFAULT_GOAL := all

all: run-soft

build:
	[ -d $(RELEASE_PATH) ] || mkdir $(RELEASE_PATH)
	DOCKER_BUILDKIT=1 $(DOCKERCMD) build . --target release --output $(RELEASE_PATH) --platform $(ENV) --progress plain

run-soft:
	$(RELEASE_PATH)/couchsport.back

run: build run-soft

clean:
	[ -d $(RELEASE_PATH) ] && rm -rf $(RELEASE_PATH)

.PHONY: build