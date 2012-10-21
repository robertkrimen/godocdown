package main

import (
	"strings"
	"bytes"
	"testing"
	"regexp"
	. "github.com/robertkrimen/terst"
)

func TestGuessImportPath(t *testing.T) {
	Terst(t)

	Is(guessImportPath(fromSlash("./example")), "github.com/robertkrimen/godocdown/godocdown/example")
	Is(guessImportPath(fromSlash("example")), "github.com/robertkrimen/godocdown/godocdown/example")
	Is(guessImportPath(fromSlash("/not/in/GOfromSlash")), "")
	Is(guessImportPath(fromSlash("in/GOfromSlash")), "github.com/robertkrimen/godocdown/godocdown/in/GOfromSlash")
	Is(guessImportPath(fromSlash("../example/example")), "github.com/robertkrimen/godocdown/example")
}

func TestFindTemplate(t *testing.T) {
	Terst(t)
	Is(findTemplate(fromSlash("../.test/godocdown.template")), fromSlash("../.test/godocdown.template/.godocdown.template"))
	Is(findTemplate(fromSlash("../.test/godocdown.tmpl")), fromSlash("../.test/godocdown.tmpl/.godocdown.tmpl"))
	Is(findTemplate(fromSlash("../.test/godocdown.md")), fromSlash("../.test/godocdown.md/.godocdown.md"))
	Is(findTemplate(fromSlash("../.test/godocdown.markdown")), fromSlash("../.test/godocdown.markdown/.godocdown.markdown"))
}

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

Also do not title something with a space at the end 

Only Title Casing Is Allowed Here

What it says on the tin above.

1word

A title with a-dash
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

Also do not title something with a space at the end 

Only Title Casing Is Allowed Here

What it says on the tin above.

# 1word

A title with a-dash
	`)

	is(synopsisHeadingTitleCase_Regexp, `
# Headline
The previous line is a single word.

a Title Is Without punctuation

	In this mode, a title can be something without punctuation

Also do not title something with a space at the end 

# Only Title Casing Is Allowed Here

What it says on the tin above.

1word

A title with a-dash
	`)

	is(synopsisHeadingTitle_Regexp, `
# Headline
The previous line is a single word.

# a Title Is Without punctuation

	In this mode, a title can be something without punctuation

Also do not title something with a space at the end 

# Only Title Casing Is Allowed Here

What it says on the tin above.

# 1word

# A title with a-dash
	`)

	is(synopsisHeadingTitleCase1Word_Regexp, `
# Headline
The previous line is a single word.

a Title Is Without punctuation

	In this mode, a title can be something without punctuation

Also do not title something with a space at the end 

# Only Title Casing Is Allowed Here

What it says on the tin above.

# 1word

A title with a-dash
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
