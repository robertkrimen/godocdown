package main

import (
	"fmt"
	"flag"
	"go/doc"
	"go/parser"
	"go/token"
	"go/printer"
	"os"
	"strings"
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"path/filepath"
)

const (
	punchCardWidth = 80
	debug = false
)

var (
	fset *token.FileSet
	synopsisHeading_Regexp = regexp.MustCompile("(?m)^([A-Za-z0-9]+)$")
	strip_Regexp = regexp.MustCompile("(?m)^\\s*// contains filtered or unexported fields\\s*\n")
	indent_Regexp = regexp.MustCompile("(?m)^([^\\n])") // Match at least one character at the start of the line

)

// Flags
var (
	signature_flag = flag.Bool("signature", false, "Add godocdown signature to the end of the documentation")
	plain_flag = flag.Bool("plain", false, "Emit standard Markdown, rather than Github Flavored Markdown (the default)")
)

type _document struct {
	name string
	pkg *doc.Package
	isCommand bool
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
	if *plain_flag {
		return indent(target + "\n", space(4))
	}
	return fmt.Sprintf("```go\n%s\n```", target)
}

func headifySynopsis(target string) string {
	return synopsisHeading_Regexp.ReplaceAllStringFunc(target, func(heading string) string {
		return "### " + heading
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

func writeConstantSection(writer io.Writer, list []*doc.Value) bool {
	empty := true
	for _, entry := range list {
		empty = false
		fmt.Fprintf(writer, "%s\n%s\n", indentCode(sourceOfNode(entry.Decl)), formatIndent(entry.Doc))
	}
	return empty
}

func writeVariableSection(writer io.Writer, list []*doc.Value) bool {
	empty := true
	for _, entry := range list {
		empty = false
		fmt.Fprintf(writer, "%s\n%s\n", indentCode(sourceOfNode(entry.Decl)), formatIndent(entry.Doc))
	}
	return empty
}

func writeFunctionSection(writer io.Writer, heading string, list []*doc.Func) bool {
	empty := true
	for _, entry := range list {
		empty = false
		receiver := " "
		if entry.Recv != "" {
			receiver = fmt.Sprintf("(%s) ", entry.Recv)
		}
		fmt.Fprintf(writer, "%s func %s%s\n\n%s\n%s\n", heading, receiver, entry.Name, indentCode(sourceOfNode(entry.Decl)), formatIndent(entry.Doc))
	}
	return empty
}

func writeTypeSection(writer io.Writer, list []*doc.Type) bool {
	empty := true
	for _, entry := range list {
		empty = false
		fmt.Fprintf(writer, "#### type %s\n\n%s\n\n%s\n", entry.Name, indentCode(sourceOfNode(entry.Decl)), formatIndent(entry.Doc))
		writeConstantSection(writer, entry.Consts)
		writeVariableSection(writer, entry.Vars)
		writeFunctionSection(writer, "####", entry.Funcs)
		writeFunctionSection(writer, "####", entry.Methods)
	}
	return empty
}

func indent(target string, indent string) string {
	return indent_Regexp.ReplaceAllString(target, indent + "$1")
}

func main() {
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		path = "."
	}

	if false {
		// Test indenting
		fmt.Printf("0/4/4\n[%s]\n",
			_formatIndent(fmt.Sprintf("%v\n%4v\n%4v\n", 1, 2, 3), space(4), space(8)))
		fmt.Printf("0/4/4\n[%s]\n",
			indent(fmt.Sprintf("%v\n%5v\n\n%5v\n", 1, 2, 3), space(4)))
		os.Exit(0)
	}

	fullPath := ""
	{
		path, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		fullPath = path
	}

	fset = token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, func(file os.FileInfo) bool {
		name := file.Name()
		if name[0] != '.' && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			return true
		}
		return false
	}, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse \"%s\": %v\n", path, err)
	}

	importLine := ""
	if read, err := ioutil.ReadFile(filepath.Join(path, ".import")); err == nil {
		importLine = strings.Split(string(read), "\n")[0]
	}

	var buffer bytes.Buffer

	// There should be only 1 package, but...
	for _, pkg := range pkgs {
		isCommand := false
		name := ""
		pkg := doc.New(pkg, ".", 0)
		switch pkg.Name {
		case "main":
			// We're probably a command, but by convention, documentation
			// should be in the documentation package:
			// http://golang.org/doc/articles/godoc_documenting_go_code.html
			continue
		case "documentation":
			// We're a command, this package/file contains the documentation
			_, name = filepath.Split(fullPath)
			isCommand = true
		default:
			name = pkg.Name
			// Just a regular package
		}

		document := &_document{
			name: name,
			pkg: pkg,
			isCommand: isCommand,
		}

		if isCommand {
			// TODO Get name from directory
		}

		// Header
		fmt.Fprintf(&buffer, "# %s\n--\n", document.name)

		if !document.isCommand {
			// Import
			if (importLine != "") {
				fmt.Fprintf(&buffer, "    import \"%s\"\n\n", importLine)
			}
		}

		// Synopsis
		fmt.Fprintf(&buffer, "%s\n", headifySynopsis(document.pkg.Doc))

		if !document.isCommand {
			// Usage
			fmt.Fprintf(&buffer, "## Usage\n\n")

			// Constant Section
			writeConstantSection(&buffer, document.pkg.Consts)

			// Variable Section
			writeVariableSection(&buffer, document.pkg.Vars)

			// Function Section
			writeFunctionSection(&buffer, "####", document.pkg.Funcs)

			// Type Section
			writeTypeSection(&buffer, document.pkg.Types)
		}
	}

	if debug {
		return
	}

	fmt.Println(strings.TrimSpace(buffer.String()))

	if *signature_flag {
		fmt.Printf("\n--\n**godocdown** http://github.com/robertkrimen/godocdown\n")
	}
}
