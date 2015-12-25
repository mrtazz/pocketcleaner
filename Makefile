#
# some housekeeping tasks
#

export GO15VENDOREXPERIMENT = 1

NAME=pocketcleaner
DESC=utility to keep your Pocket list small and manageable. Archives everything in your list except for the newest n items.
PREFIX ?= /usr/local
VERSION=$(shell git describe --tags --always --dirty)
GOVERSION = $(shell go version)
BUILDTIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDER = $(shell echo "`git config user.name` <`git config user.email`>")
PKG_RELEASE ?= 1
PROJECT_URL="https://github.com/mrtazz/$(NAME)"
SOURCES=cmd/pocketcleaner/main.go pocketcleaner.go
LDFLAGS=-X 'main.version=$(VERSION)' \
				-X 'main.buildTime=$(BUILDTIME)'\
				-X 'main.builder=$(BUILDER)'\
				-X 'main.goversion=$(GOVERSION)'
TARGETS=$(PREFIX)/bin/$(NAME) $(PREFIX)/share/man/man1/$(NAME).1

$(NAME): $(SOURCES)
	go build -ldflags "$(LDFLAGS)" -o $@ $<

$(PREFIX)/bin:
	install -m 755 -d $@

$(PREFIX)/bin/$(NAME): $(NAME) $(PREFIX)/bin
	install -m 755 $< $@

$(NAME).1: $(NAME).1.txt
	txt2man -t "$(NAME)" -s 1 -v "User Manual" $< > $@

$(PREFIX)/share/man/man1:
	install -m 755 -d $@

$(PREFIX)/share/man/man1/$(NAME).1: $(NAME).1 $(PREFIX)/share/man/man1
	install -m 755 $< $@

.PHONY: test rpm deb local-install packages coverage vet

install: $(TARGETS)

local-install:
	$(MAKE) install PREFIX=usr

packages: local-install rpm deb

deploy-packages: packages
	package_cloud push mrtazz/$(NAME)/el/7 *.rpm
	package_cloud push mrtazz/$(NAME)/debian/wheezy *.deb
	package_cloud push mrtazz/$(NAME)/ubuntu/trusty *.deb


rpm: $(SOURCES)
	  fpm -t rpm -s dir \
    --name $(NAME) \
    --version $(VERSION) \
		--description "$(DESC)" \
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
		--description "$(DESC)" \
    --iteration $(PKG_RELEASE) \
    --epoch 1 \
    --license MIT \
    --maintainer "Daniel Schauenberg <d@unwiredcouch.com>" \
    --url $(PROJECT_URL) \
    --vendor mrtazz \
    usr

GAUGES_CODE=5678b0854b2ffa74ed002b8e

jekyll:
	install -d ./docs
	echo "gaugesid: $(GAUGES_CODE)" > docs/_config.yml
	echo "projecturl: $(PROJECT_URL)" >> docs/_config.yml
	echo "basesite: http://www.unwiredcouch.com" >> docs/_config.yml
	echo "markdown: redcarpet" >> docs/_config.yml
	echo "---" > docs/index.md
	echo "layout: project" >> docs/index.md
	echo "title: $(NAME)" >> docs/index.md
	echo "---" >> docs/index.md
	cat README.md >> docs/index.md

docs: jekyll

clean-docs:
	rm -rf ./docs

deploy-docs: docs
	@cd docs && git init && git remote add upstream "https://${GH_TOKEN}@github.com/mrtazz/$(NAME).git" && \
	git submodule add https://github.com/mrtazz/jekyll-layouts.git ./_layouts && \
	git submodule update --init && \
	git fetch upstream && git reset upstream/gh-pages && \
	git config user.name 'Daniel Schauenberg' && git config user.email d@unwiredcouch.com && \
	touch . && git add -A . && \
	git commit -m "rebuild pages at $(VERSION)" && \
	git push -q upstream HEAD:gh-pages

clean: clean-docs
	rm -rf ./usr
	rm $(NAME)
	rm $(NAME).1

test:
	@go test -v .

vet:
	@go tool vet .

coverage:
	@-go test -v -coverprofile=cover.out github.com/mrtazz/$(NAME)
	@-go tool cover -html=cover.out -o cover.html

benchmark:
	@echo "Running tests..."
	@go test -bench=. ${NAME}

format:
	@go fmt .

govendor:
	    go get -u github.com/kardianos/govendor
