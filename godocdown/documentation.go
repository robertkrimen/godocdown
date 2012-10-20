/*
Command godocdown extracts and generates Go documentation in a GitHub-friendly Markdown format.

	$ go get github.com/robertkrimen/godocdown/godocdown

	$ godocdown /path/to/package > README.markdown

	# Generate documentation for the package/command in the current directory
	$ godocdown > README.markdown

	# Generate standard Markdown
	$ godocdown -plain . 

This program is targeted at providing nice-looking documentation for GitHub. With this in
mind, it generates GitHub Flavored Markdown (http://github.github.com/github-flavored-markdown/) by
default. This can be changed with the use of the "plain" flag to generate standard Markdown.

Installation

	go get github.com/robertkrimen/godocdown/godocdown

Example

http://github.com/robertkrimen/godocdown/blob/master/example.markdown

Usage

The following options are accepted:

	-heading="TitleCase1Word"
	// Heading detection method: 1Word, TitleCase, Title, TitleCase1Word, ""
	// For each line of the package declaration, godocdown attempts to detect if
	// a heading is present via a pattern match. If a heading is detected,
	// it prefixes the line with a Markdown heading indicator (typically "###").

		1Word: Only a single word on the entire line
			[A-Za-z0-9_-]+

		TitleCase: A line where each word has the first letter capitalized
			([A-Z][A-Za-z0-9_-]\s*)+

		Title: A line without punctuation (e.g. a period at the end)
			([A-Za-z0-9_-]\s*)+

		TitleCase1Word: The line matches either the TitleCase or 1Word pattern

	-no-template=false
	// Disable template processing

	-plain=false
	// Emit standard Markdown, rather than Github Flavored Markdown

Templating

In addition to Markdown rendering, godocdown provides templating via text/template (http://golang.org/pkg/text/template/)
for further customization. By putting a file named ".godocdown.template" in the same directory as your
package/command, godocdown will know to use the file as a template. The following names are also accepted,
with the first one encountered being used:

	# text/template
	.godocdown.markdown
	.godocdown.md
	.godocdown.template
	.godocdown.tmpl

In addition to the standard template functionality, the starting data argument has the following interface:

	.Emit
	// A method for emitting the standard documentation (what godocdown would emit without a template)

	.EmitHeader
	// A method for emitting the package name and an import line (if one is present/needed)

	.EmitSynopsis
	// A method for emitting the package declaration

	.EmitUsage
	// A method for emitting package usage, which includes a constants section, a variables section,
	// a functions section, and a types section. In addition, each type may have its own constant,
	// variable, and/or function/method listing.

	.IsCommand
	// A boolean indicating whether the given package is a command or a plain package

	.Name
	// A string containing the name of the package/command

*/
package documentation
