/*
Command godocdown extracts and generates Go documentation in a GitHub-friendly Markdown format.

This program is targeted at providing nice-looking documentation for GitHub. With this in
mind, it generates GitHub Flavored Markdown (http://github.github.com/github-flavored-markdown/) by
default. This can be changed with the use of the "plain" flag to generate standard Markdown.

	$ go get github.com/robertkrimen/godocdown/godocdown

	$ godocdown /path/to/package > README.markdown

	# Generate documentation for the package/command in the current directory
	$ godocdown > README.markdown

	# Generate standard Markdown
	$ godocdown -plain . 

Installation

	go get github.com/robertkrimen/godocdown/godocdown

Example

http://github.com/robertkrimen/godocdown/blob/master/example.markdown

Usage

The following options are accepted:

	-heading="TitleCase1Word"
	// Heading detection method: 1Word, TitleCase, Title, TitleCase1Word, ""

	-no-template=false
	// Disable template processing

	-plain=false
	// Emit standard Markdown, rather than Github Flavored Markdown (the default)

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
	// A method emitting all of the rendered documentation

	.EmitHeader
	// A method emitting the package/command name and an import line (if one is present/needed)

	.EmitSynopsis
	// A method emitting the package/command declaration

	.EmitUsage
	// A method emitting package usage, include a constant section, a variable section,
	// the function section, and a type section (and each type having its own constant, variable,
	// and function/method listing)

	.IsCommand
	// A boolean indicating whether the given package is a command or a plain package

	.Name
	// A string containing the name of the package/command

*/
package documentation
