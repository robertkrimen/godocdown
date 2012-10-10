.PHONY: test build test-example release

TEST := ./example

test:
	go run main.go $(TEST)

build:
	go build -o godocdown

test-example: build
	./godocdown example > test/README.markdown
	cd test && git commit -m 'WIP' * && git push

release:
	./godocdown $(HOME)/go/src/pkg/strings > example.markdown
