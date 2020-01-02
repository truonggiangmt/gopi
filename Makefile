# Go parameters
GO=go

# App parameters
GOPI=github.com/djthorpe/gopi/v2/config
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

all:
	@echo "Synax: make linux|darwin|rpi"

# Build for different platforms
linux: TAGS = -tags linux
linux: test install

rpi: TAGS = -tags rpi
rpi: PKG_CONFIG_PATH = /opt/vc/lib/pkgconfig
rpi: test install

darwin: TAGS = -tags darwin
darwin: test install

# Build rules
test: 
	PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" $(GO) test $(TAGS) -v ./...

install:
	PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" $(GO) install $(TAGS) ${GOFLAGS} ./cmd/...

clean: 
	$(GO) clean

