.PHONY: test build test-example release run-example run-help install

RUN := ./run
TEST := ./example

export TERST_BASE=$(PWD)

test:
	go test -i ./godocdown &&  go test -v ./godocdown

build:
	cd godocdown && go build -o ../$(RUN)

test-example: build
	$(RUN) -signature example > test/README.markdown
	cd test && git commit -m 'WIP' * && git push

release: test build 
	$(RUN) $(HOME)/go/src/pkg/strings > example.markdown
	$(RUN) -signature ./godocdown > README.markdown

run-help:
	cd godocdown && go run main.go render.go -help

run-example:
	cd godocdown && go run main.go render.go --signature ../$(TEST)

install:
	go install ./godocdown
