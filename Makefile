GO ?= go
GOPATH := $(CURDIR):$(CURDIR)/Godeps/_workspace:$(GOPATH)

all: build

build:
	$(GO) build -o httpception httpception/httpception

clean:
	rm httpception
