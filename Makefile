.PHONY: build test man install

build:
	go build -o lookit ./cmd/lookit

test:
	go test ./...

man: build
	./lookit gen-man ./man/man1

install: build man
	install -m 755 lookit /usr/local/bin/
	install -d /usr/local/share/man/man1/
	install -m 644 man/man1/*.1 /usr/local/share/man/man1/
