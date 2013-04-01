.PHONY: build test test-example install release

build: test
	go build -v

test:
	go test -i
	go test -v

test-example: build
	./godocdown -signature example > test/README.markdown
	#cd test && git commit -m 'WIP' * && git push

install:
	go install -v

release:
	$(MAKE) -C .. $@

