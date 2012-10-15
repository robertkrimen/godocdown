package main

import (
	"strings"
	"bytes"
	"testing"
	"regexp"
	. "github.com/robertkrimen/terst"
)

func TestIndent(t *testing.T) {
	Terst(t)

	Is(indent("1\n  2\n\n  3\n  4\n", "  "), "  1\n    2\n\n    3\n    4\n")
}

func TestHeadlineSynopsis(t *testing.T) {
	Terst(t)

	synopsis := `
Headline
The previous line is a single word.

a Title Is Without punctuation

	In this mode, a title can be something without punctuation

Only Title Casing Is Allowed Here

What it says on the tin above.
	`
	is := func(scanner *regexp.Regexp, want string){
		have := headlineSynopsis(synopsis, "#", scanner)
		Is(strings.TrimSpace(have), strings.TrimSpace(want))
	}

	is(synopsisHeading1Word_Regexp, `
# Headline
The previous line is a single word.

a Title Is Without punctuation

	In this mode, a title can be something without punctuation

Only Title Casing Is Allowed Here

What it says on the tin above.
	`)

	is(synopsisHeadingTitleCase_Regexp, `
# Headline
The previous line is a single word.

a Title Is Without punctuation

	In this mode, a title can be something without punctuation

# Only Title Casing Is Allowed Here

What it says on the tin above.
	`)

	is(synopsisHeadingTitle_Regexp, `
# Headline
The previous line is a single word.

# a Title Is Without punctuation

	In this mode, a title can be something without punctuation

# Only Title Casing Is Allowed Here

What it says on the tin above.
	`)
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

	renderSignatureTo(buffer)
	Is(buffer.String(), "\n\n--\n**godocdown** http://github.com/robertkrimen/godocdown\n")
}
