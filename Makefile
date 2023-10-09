GO=go
GOBUILD=$(GO) build

.PHONY: build

build:
	$(GOBUILD) -o ./build/ ./cmd/...

.PHONY: clean

clean:
	-rm -r ./build
