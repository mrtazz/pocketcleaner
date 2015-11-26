#
# some housekeeping tasks
#

export GO15VENDOREXPERIMENT = 1

NAME=pocketcleaner
PREFIX ?= /usr/local
VERSION=$(shell git describe --tags --always --dirty)
GOVERSION = $(shell go version)
BUILDTIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDER = $(shell echo "`git config user.name` <`git config user.email`>")
PKG_RELEASE ?= 1
PROJECT_URL="https://github.com/mrtazz/$(NAME)"
SOURCES=cmd/pocketcleaner/main.go
LDFLAGS=-X 'main.version=$(VERSION)' \
				-X 'main.buildTime=$(BUILDTIME)'\
				-X 'main.builder=$(BUILDER)'\
				-X 'main.goversion=$(GOVERSION)'

$(NAME): $(SOURCES)
	go build -ldflags "$(LDFLAGS)" -o $@ $<

$(PREFIX)/bin:
	install -m 755 -d $@

$(PREFIX)/bin/$(NAME): $(NAME) $(PREFIX)/bin
	install -m 755 $< $@

.PHONY: test rpm deb local-install packages

local-install:
	$(MAKE) install PREFIX=usr

packages: local-install rpm deb

rpm: $(SOURCES)
	  fpm -t rpm -s dir \
    --name $(NAME) \
    --version $(VERSION) \
    --iteration $(PKG_RELEASE) \
    --epoch 1 \
    --license MIT \
    --maintainer "Daniel Schauenberg <d@unwiredcouch.com>" \
    --url $(PROJECT_URL) \
    --vendor mrtazz \
    usr

deb: $(SOURCES)
	  fpm -t deb -s dir \
    --name $(NAME) \
    --version $(VERSION) \
    --iteration $(PKG_RELEASE) \
    --epoch 1 \
    --license MIT \
    --maintainer "Daniel Schauenberg <d@unwiredcouch.com>" \
    --url $(PROJECT_URL) \
    --vendor mrtazz \
    usr

test:
	@go test .

coverage:
	@go test -coverprofile=cover.out github.com/mrtazz/$(NAME)
	@go tool cover -html=cover.out -o cover.html

benchmark:
	@echo "Running tests..."
	@go test -bench=. ${NAME}

format:
	@go fmt .

govendor:
	    go get -u github.com/kardianos/govendor
