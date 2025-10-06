PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin

.PHONY: build install clean test

build:
	go build -o fling

install: build
	install -d $(BINDIR)
	install -m 755 fling $(BINDIR)/fling

clean:
	rm -f fling

test: build
	@echo "Built successfully. Testing PATH scanning..."
	@echo "Found $(shell go run . 2>/dev/null | wc -l) applications in PATH"
	@echo "To test interactively, run: ./fling"