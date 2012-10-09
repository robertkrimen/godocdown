.PHONY: test build style

TEST := ./example

test:
	go run main.go $(TEST)

build:
	go build -o godocdown

style: build
	./godocdown example > test/README.markdown
	cd test && git commit -m 'WIP' * && git push
