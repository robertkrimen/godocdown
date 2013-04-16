.PHONY: build test test-example install release

build: test
	go build

test:
	go test -i
	go test

test-example: build
	./godocdown -signature example > test/README.markdown
	#cd test && git commit -m 'WIP' * && git push

install:
	go install

release:
	$(MAKE) -C .. $@
