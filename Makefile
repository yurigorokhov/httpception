GO ?= godep go
GOPATH := $(CURDIR)/Godeps/_workspace:$(GOPATH)

all: build

build:
	$(GO) build -o httpception httpception

clean:
	rm httpception
