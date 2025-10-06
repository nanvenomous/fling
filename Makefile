PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin

.PHONY: build install clean test

build:
	go build -o fling

install: build
	install -d $(BINDIR)
	install -m 755 fling $(BINDIR)/fling
	install -m 755 apps.sh $(BINDIR)/fling-apps.sh

clean:
	rm -f fling

test: build
	./apps.sh all | head -5
	@echo "Built successfully. To test interactively, run: echo 'firefox\nchromium\nvim' | fzf"