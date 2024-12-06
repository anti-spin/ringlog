# Load variables from .env if it exists
ifneq (,$(wildcard .env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

BINARY_NAME=ringlog
GO_MODULE=github.com/anti-spin/ringlog
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

# Get the latest git tag or commit hash
VERSION := $(shell git describe --tags --always --dirty)

BUILD_DIR=build
DEB_DIR=$(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_amd64

# Default target
all: clean package upload

# Build the binary with static linking
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -ldflags="-X $(GO_MODULE)/cmd.version=$(VERSION)" -o $(BINARY_NAME)
	@echo "Build completed: $(BINARY_NAME)"

# Clean the binary
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	@echo "Clean completed."

package-prepare: build
	mkdir -p $(DEB_DIR)/DEBIAN
	install -m 755 -D $(BINARY_NAME) $(DEB_DIR)/usr/local/bin/$(BINARY_NAME)
	install -m 644 -D debian/$(BINARY_NAME).1 $(DEB_DIR)/usr/share/man/man1/$(BINARY_NAME).1
	sed "s/@VERSION@/$(VERSION)/g" < debian/control.template > $(DEB_DIR)/DEBIAN/control
	if [ -f debian/postinst ]; then install -m 755 -D debian/postinst $(DEB_DIR)/DEBIAN/postinst; fi
	if [ -f debian/postrm ]; then install -m 755 -D debian/postrm $(DEB_DIR)/DEBIAN/postrm; fi

package: package-prepare
	@echo "Packaging $(BINARY_NAME)..."
	dpkg-deb --build $(DEB_DIR)
	@echo "Created $(BINARY_NAME)_$(VERSION)_amd64.deb"

upload:
	@echo "Checking for DEB_REGISTRY..."
	@if [ -n "$(DEB_REGISTRY)" ]; then \
		echo "Uploading to registry $(DEB_REGISTRY)"; \
		./upload-deb.sh $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_amd64.deb; \
	else \
		echo "DEB_REGISTRY not set, skipping upload"; \
	fi
