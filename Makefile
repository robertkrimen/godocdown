.PHONY: test build test-example release

RUN := ./run
TEST := ./example

export TERST_BASE=$(PWD)

test:
	go test -i ./godocdown &&  go test ./godocdown
	go run godocdown/main.go $(TEST)

build:
	cd godocdown && go build -o ../$(RUN)

test-example: build
	$(RUN) --signature example > test/README.markdown
	cd test && git commit -m 'WIP' * && git push

release: build
	$(RUN) $(HOME)/go/src/pkg/strings > example.markdown
