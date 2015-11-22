#
# some housekeeping tasks
#

export GO15VENDOREXPERIMENT = 1

NAME=pocket-cleaner
PREFIX ?= /usr/local
VERSION=$(shell git describe --abbrev=0)
GOVERSION = $(shell go version)
BUILDTIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDER = $(shell echo "`git config user.name` <`git config user.email`>")
PKG_RELEASE ?= 1
PROJECT_URL="https://github.com/mrtazz/$(NAME)"
SOURCES=main.go
LDFLAGS=-X 'main.version=$(VERSION)' \
				-X 'main.buildTime=$(BUILDTIME)'\
				-X 'main.builder=$(BUILDER)'\
				-X 'main.goversion=$(GOVERSION)'

$(NAME): main.go
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

benchmark:
	@echo "Running tests..."
	@go test -bench=. ${NAME}

format:
	@go fmt .

govendor:
	    go get -u github.com/kardianos/govendor
