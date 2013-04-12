/*
Command godocdown generates Go documentation in a GitHub-friendly Markdown format.

    $ go get github.com/robertkrimen/godocdown/godocdown

    $ godocdown /path/to/package > README.markdown

    # Generate documentation for the package/command in the current directory
    $ godocdown > README.markdown

    # Generate standard Markdown
    $ godocdown -plain . 

This program is targeted at providing nice-looking documentation for GitHub. With this in
mind, it generates GitHub Flavored Markdown (http://github.github.com/github-flavored-markdown/) by
default. This can be changed with the use of the "plain" flag to generate standard Markdown.

Install

    go get github.com/robertkrimen/godocdown/godocdown

Example

http://github.com/robertkrimen/godocdown/blob/master/example.markdown

Usage

    -template="": The template file to use

    -no-template=false
        Disable template processing

    -plain=false
        Emit standard Markdown, rather than Github Flavored Markdown

    -heading="TitleCase1Word"
        Heading detection method: 1Word, TitleCase, Title, TitleCase1Word, ""
        For each line of the package declaration, godocdown attempts to detect if
        a heading is present via a pattern match. If a heading is detected,
        it prefixes the line with a Markdown heading indicator (typically "###").

        1Word: Only a single word on the entire line
            [A-Za-z0-9_-]+

        TitleCase: A line where each word has the first letter capitalized
            ([A-Z][A-Za-z0-9_-]\s*)+

        Title: A line without punctuation (e.g. a period at the end)
            ([A-Za-z0-9_-]\s*)+

        TitleCase1Word: The line matches either the TitleCase or 1Word pattern

Templating

In addition to Markdown rendering, godocdown provides templating via text/template (http://golang.org/pkg/text/template/)
for further customization. By putting a file named ".godocdown.template" (or one from the list below) in the same directory as your
package/command, godocdown will know to use the file as a template.

    # text/template
    .godocdown.markdown
    .godocdown.md
    .godocdown.template
    .godocdown.tmpl

A template file can also be specified with the "-template" parameter

Along with the standard template functionality, the starting data argument has the following interface:

    {{ .Emit }}
    // Emit the standard documentation (what godocdown would emit without a template)

    {{ .EmitHeader }}
    // Emit the package name and an import line (if one is present/needed)

    {{ .EmitSynopsis }}
    // Emit the package declaration

    {{ .EmitUsage }}
    // Emit package usage, which includes a constants section, a variables section,
    // a functions section, and a types section. In addition, each type may have its own constant,
    // variable, and/or function/method listing.

    {{ if .IsCommand  }} ... {{ end }}
    // A boolean indicating whether the given package is a command or a plain package

    {{ .Name }}
    // The name of the package/command (string)

    {{ .ImportPath }}
    // The import path for the package (string)
    // (This field will be the empty string if godocdown is unable to guess it)

*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	tmplate "text/template"
	tme "time"
)

const (
	punchCardWidth = 80
	debug          = false
)

var (
	flag_signature  = flag.Bool("signature", false, "Add godocdown signature to the end of the documentation")
	flag_plain      = flag.Bool("plain", false, "Emit standard Markdown, rather than Github Flavored Markdown (the default)")
	flag_heading    = flag.String("heading", "TitleCase1Word", "Heading detection method: 1Word, TitleCase, Title, TitleCase1Word, \"\"")
	flag_template   = flag.String("template", "", "The template file to use")
	flag_noTemplate = flag.Bool("no-template", false, "Disable template processing")
)

var (
	fset *token.FileSet

	synopsisHeading1Word_Regexp          = regexp.MustCompile("(?m)^([A-Za-z0-9_-]+)$")
	synopsisHeadingTitleCase_Regexp      = regexp.MustCompile("(?m)^((?:[A-Z][A-Za-z0-9_-]*)(?:[ \t]+[A-Z][A-Za-z0-9_-]*)*)$")
	synopsisHeadingTitle_Regexp          = regexp.MustCompile("(?m)^((?:[A-Za-z0-9_-]+)(?:[ \t]+[A-Za-z0-9_-]+)*)$")
	synopsisHeadingTitleCase1Word_Regexp = regexp.MustCompile("(?m)^((?:[A-Za-z0-9_-]+)|(?:(?:[A-Z][A-Za-z0-9_-]*)(?:[ \t]+[A-Z][A-Za-z0-9_-]*)*))$")

	strip_Regexp           = regexp.MustCompile("(?m)^\\s*// contains filtered or unexported fields\\s*\n")
	indent_Regexp          = regexp.MustCompile("(?m)^([^\\n])") // Match at least one character at the start of the line
	synopsisHeading_Regexp = synopsisHeading1Word_Regexp
)

var DefaultStyle = Style{
	IncludeImport: true,

	SynopsisHeader:  "###",
	SynopsisHeading: synopsisHeadingTitleCase1Word_Regexp,

	UsageHeader: "## Usage\n",

	ConstantHeader:     "####",
	VariableHeader:     "####",
	FunctionHeader:     "####",
	TypeHeader:         "####",
	TypeFunctionHeader: "####",

	IncludeSignature: false,
}
var RenderStyle = DefaultStyle

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	executable, err := os.Stat(os.Args[0])
	if err != nil {
		return
	}
	time := executable.ModTime()
	since := tme.Since(time)
	fmt.Fprintf(os.Stderr, "---\n%s (%.2f)\n", time.Format("2006-01-02 15:04 MST"), since.Minutes())
}

func init() {
	flag.Usage = usage
}

type Style struct {
	IncludeImport bool

	SynopsisHeader  string
	SynopsisHeading *regexp.Regexp

	UsageHeader string

	ConstantHeader     string
	VariableHeader     string
	FunctionHeader     string
	TypeHeader         string
	TypeFunctionHeader string

	IncludeSignature bool
}

type _document struct {
	Name       string
	pkg        *doc.Package
	buildPkg   *build.Package
	IsCommand  bool
	ImportPath string
}

func _formatIndent(target, indent, preIndent string) string {
	var buffer bytes.Buffer
	doc.ToText(&buffer, target, indent, preIndent, punchCardWidth-2*len(indent))
	return buffer.String()
}

func space(width int) string {
	return strings.Repeat(" ", width)
}

func formatIndent(target string) string {
	return _formatIndent(target, space(0), space(4))
}

func indentCode(target string) string {
	if *flag_plain {
		return indent(target+"\n", space(4))
	}
	return fmt.Sprintf("```go\n%s\n```", target)
}

func headifySynopsis(target string) string {
	detect := RenderStyle.SynopsisHeading
	if detect == nil {
		return target
	}
	return detect.ReplaceAllStringFunc(target, func(heading string) string {
		return fmt.Sprintf("%s %s", RenderStyle.SynopsisHeader, heading)
	})
}

func headlineSynopsis(synopsis, header string, scanner *regexp.Regexp) string {
	return scanner.ReplaceAllStringFunc(synopsis, func(headline string) string {
		return fmt.Sprintf("%s %s", header, headline)
	})
}

func sourceOfNode(target interface{}) string {
	var buffer bytes.Buffer
	mode := printer.TabIndent | printer.UseSpaces
	err := (&printer.Config{Mode: mode, Tabwidth: 4}).Fprint(&buffer, fset, target)
	if err != nil {
		return ""
	}
	return strip_Regexp.ReplaceAllString(buffer.String(), "")
}

func indent(target string, indent string) string {
	return indent_Regexp.ReplaceAllString(target, indent+"$1")
}

func trimSpace(buffer *bytes.Buffer) {
	tmp := bytes.TrimSpace(buffer.Bytes())
	buffer.Reset()
	buffer.Write(tmp)
}

func fromSlash(path string) string {
	return filepath.FromSlash(path)
}

/*
    This is how godoc does it:

	// Determine paths.
	//
	// If we are passed an operating system path like . or ./foo or /foo/bar or c:\mysrc,
	// we need to map that path somewhere in the fs name space so that routines
	// like getPageInfo will see it.  We use the arbitrarily-chosen virtual path "/target"
	// for this.  That is, if we get passed a directory like the above, we map that
	// directory so that getPageInfo sees it as /target.
	const target = "/target"
	const cmdPrefix = "cmd/"
	path := flag.Arg(0)
	var forceCmd bool
	var abspath, relpath string
	if filepath.IsAbs(path) {
		fs.Bind(target, OS(path), "/", bindReplace)
		abspath = target
	} else if build.IsLocalImport(path) {
		cwd, _ := os.Getwd() // ignore errors
		path = filepath.Join(cwd, path)
		fs.Bind(target, OS(path), "/", bindReplace)
		abspath = target
	} else if strings.HasPrefix(path, cmdPrefix) {
		path = path[len(cmdPrefix):]
		forceCmd = true
	} else if bp, _ := build.Import(path, "", build.FindOnly); bp.Dir != "" && bp.ImportPath != "" {
		fs.Bind(target, OS(bp.Dir), "/", bindReplace)
		abspath = target
		relpath = bp.ImportPath
	} else {
		abspath = pathpkg.Join(pkgHandler.fsRoot, path)
	}
	if relpath == "" {
		relpath = abspath
	}
*/
func buildImport(target string) (*build.Package, error) {
	if filepath.IsAbs(target) {
		return build.Default.ImportDir(target, build.FindOnly)
	} else if build.IsLocalImport(target) {
		base, _ := os.Getwd()
		path := filepath.Join(base, target)
		return build.Default.ImportDir(path, build.FindOnly)
	} else if pkg, _ := build.Default.Import(target, "", build.FindOnly); pkg.Dir != "" && pkg.ImportPath != "" {
		return pkg, nil
	}
	path, _ := filepath.Abs(target) // Even if there is an error, still try?
	return build.Default.ImportDir(path, build.FindOnly)
}

func guessImportPath(target string) (string, error) {
	buildPkg, err := buildImport(target)
	if err != nil {
		return "", err
	}
	if buildPkg.SrcRoot == "" {
		return "", nil
	}
	return buildPkg.ImportPath, nil
}

func loadDocument(target string) (*_document, error) {

	buildPkg, err := buildImport(target)
	if err != nil {
		return nil, err
	}
	if buildPkg.Dir == "" {
		return nil, fmt.Errorf("Could not find package \"%s\"", target)
	}

	path := buildPkg.Dir

	fset = token.NewFileSet()
	pkgSet, err := parser.ParseDir(fset, path, func(file os.FileInfo) bool {
		name := file.Name()
		if name[0] != '.' && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			return true
		}
		return false
	}, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("Could not parse \"%s\": %v", path, err)
	}

	importPath := ""
	if read, err := ioutil.ReadFile(filepath.Join(path, ".godocdown.import")); err == nil {
		importPath = strings.TrimSpace(strings.Split(string(read), "\n")[0])
	} else {
		importPath = buildPkg.ImportPath
	}

	{
		isCommand := false
		name := ""
		var pkg *doc.Package

		// Choose the best package for documentation. Either
		// documentation, main, or whatever the package is.
		for _, parsePkg := range pkgSet {
			tmpPkg := doc.New(parsePkg, ".", 0)
			switch tmpPkg.Name {
			case "main":
				if isCommand {
					// We've already seen "pacakge documentation",
					// so favor that over main.
					continue
				}
				fallthrough
			case "documentation":
				// We're a command, this package/file contains the documentation
				// path is used to get the containing directory in the case of
				// command documentation
				path, err := filepath.Abs(path)
				if err != nil {
					panic(err)
				}
				_, name = filepath.Split(path)
				isCommand = true
				pkg = tmpPkg
			default:
				// Just a regular package
				name = tmpPkg.Name
				pkg = tmpPkg
			}
		}

		if pkg != nil {
			return &_document{
				Name:       name,
				pkg:        pkg,
				buildPkg:   buildPkg,
				IsCommand:  isCommand,
				ImportPath: importPath,
			}, nil
		}
	}

	return nil, nil
}

func emitString(fn func(*bytes.Buffer)) string {
	var buffer bytes.Buffer
	fn(&buffer)
	trimSpace(&buffer)
	return buffer.String()
}

// Emit
func (self *_document) Emit() string {
	return emitString(func(buffer *bytes.Buffer) {
		self.EmitTo(buffer)
	})
}

func (self *_document) EmitTo(buffer *bytes.Buffer) {

	// Header
	self.EmitHeaderTo(buffer)

	// Synopsis
	self.EmitSynopsisTo(buffer)

	// Usage
	if !self.IsCommand {
		self.EmitUsageTo(buffer)
	}

	trimSpace(buffer)
}

// Signature
func (self *_document) EmitSignature() string {
	return emitString(func(buffer *bytes.Buffer) {
		self.EmitSignatureTo(buffer)
	})
}

func (self *_document) EmitSignatureTo(buffer *bytes.Buffer) {

	renderSignatureTo(buffer)

	trimSpace(buffer)
}

// Header
func (self *_document) EmitHeader() string {
	return emitString(func(buffer *bytes.Buffer) {
		self.EmitHeaderTo(buffer)
	})
}

func (self *_document) EmitHeaderTo(buffer *bytes.Buffer) {
	renderHeaderTo(buffer, self)
}

// Synopsis
func (self *_document) EmitSynopsis() string {
	return emitString(func(buffer *bytes.Buffer) {
		self.EmitSynopsisTo(buffer)
	})
}

func (self *_document) EmitSynopsisTo(buffer *bytes.Buffer) {
	renderSynopsisTo(buffer, self)
}

// Usage
func (self *_document) EmitUsage() string {
	return emitString(func(buffer *bytes.Buffer) {
		self.EmitUsageTo(buffer)
	})
}

func (self *_document) EmitUsageTo(buffer *bytes.Buffer) {
	renderUsageTo(buffer, self)
}

var templateNameList = strings.Fields(`
	.godocdown.markdown
	.godocdown.md
	.godocdown.template
	.godocdown.tmpl
`)

func findTemplate(path string) string {

	for _, templateName := range templateNameList {
		templatePath := filepath.Join(path, templateName)
		_, err := os.Stat(templatePath)
		if err != nil {
			if os.IsExist(err) {
				continue
			}
			continue // Other error reporting?
		}
		return templatePath
	}
	return "" // Nothing found
}

func loadTemplate(document *_document) *tmplate.Template {
	if *flag_noTemplate {
		return nil
	}

	templatePath := *flag_template
	if templatePath == "" {
		templatePath = findTemplate(document.buildPkg.Dir)
	}

	if templatePath == "" {
		return nil
	}

	template := tmplate.New("").Funcs(tmplate.FuncMap{})
	template, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template \"%s\": %v", templatePath, err)
		os.Exit(1)
	}
	return template
}

func main() {
	flag.Parse()
	target := flag.Arg(0)
	fallbackUsage := false
	if target == "" {
		fallbackUsage = true
		target = "."
	}

	RenderStyle.IncludeSignature = *flag_signature

	switch *flag_heading {
	case "1Word":
		RenderStyle.SynopsisHeading = synopsisHeading1Word_Regexp
	case "TitleCase":
		RenderStyle.SynopsisHeading = synopsisHeadingTitleCase_Regexp
	case "Title":
		RenderStyle.SynopsisHeading = synopsisHeadingTitle_Regexp
	case "TitleCase1Word":
		RenderStyle.SynopsisHeading = synopsisHeadingTitleCase1Word_Regexp
	case "", "-":
		RenderStyle.SynopsisHeading = nil
	}

	document, err := loadDocument(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	if document == nil {
		// Nothing found.
		if fallbackUsage {
			usage()
			os.Exit(2)
		} else {
			fmt.Fprintf(os.Stderr, "Could not find package: %s\n", target)
			os.Exit(1)
		}
	}

	template := loadTemplate(document)

	var buffer bytes.Buffer
	if template == nil {
		document.EmitTo(&buffer)
		document.EmitSignatureTo(&buffer)
	} else {
		err := template.Templates()[0].Execute(&buffer, document)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running template: %v", err)
			os.Exit(1)
		}
		document.EmitSignatureTo(&buffer)
	}

	if debug {
		// Skip printing if we're debugging
		return
	}

	documentation := buffer.String()
	documentation = strings.TrimSpace(documentation)
	fmt.Println(documentation)
}
