.PHONY: test build test-example release

TEST := ./example

test:
	go run godocdown/main.go $(TEST)

build:
	cd godocdown && go build -o ../.godocdown

test-example: build
	./.godocdown --signature example > test/README.markdown
	cd test && git commit -m 'WIP' * && git push

release: build
	./.godocdown $(HOME)/go/src/pkg/strings > example.markdown
