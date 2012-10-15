package main

import (
	"strings"
	"bytes"
	"testing"
	. "github.com/robertkrimen/terst"
)

func TestIndent(t *testing.T) {
	Terst(t)

	Is(indent("1\n  2\n\n  3\n  4\n", "  "), "  1\n    2\n\n    3\n    4\n")
}

func Test(t *testing.T) {
	Terst(t)

	document, err := loadDocument("../example")
	if err != nil {
		Is(err.Error(), "")
		return
	}
	if document == nil {
		Is("200", "404") // Heh
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	is := func(want string){
		Is(strings.TrimSpace(buffer.String()), strings.TrimSpace(want))
		buffer.Reset()
	}

	renderHeaderTo(buffer, document)
	is(`
# example
--
    import "github.com/robertkrimen/godocdown/example"
	`)

	RenderStyle.IncludeImport = false
	renderHeaderTo(buffer, document)
	is(`
# example
--
	`)

	renderSynopsisTo(buffer, document)
	is(`
Package example is an example package with documentation

	// Here is some code
	func example() {
		abc := 1 + 1
	}()

### Installation

	# This is how to install it:
	$ curl http://example.com
	$ tar xf example.tar.gz -C .
	$ ./example &
	`)

	RenderStyle.IncludeSignature = true
	renderSignatureTo(buffer)
	is(`
--
**godocdown** http://github.com/robertkrimen/godocdown
	`)
}
