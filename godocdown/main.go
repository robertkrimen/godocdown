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

// Flags
var (
	signature_flag = flag.Bool("signature", false, "Add godocdown signature to the end of the documentation")
	plain_flag = flag.Bool("plain", false, "Emit standard Markdown, rather than Github Flavored Markdown (the default)")
)

var (
	fset *token.FileSet
	synopsisHeading_Regexp = regexp.MustCompile("(?m)^([A-Za-z0-9]+)$")
	strip_Regexp = regexp.MustCompile("(?m)^\\s*// contains filtered or unexported fields\\s*\n")
	indent_Regexp = regexp.MustCompile("(?m)^([^\\n])") // Match at least one character at the start of the line

)

var Style = struct {
	IncludeImport bool

	SynopsisHeader string
	HeadifySynopsis bool

	ConstantHeader string
	VariableHeader string
	FunctionHeader string
	TypeHeader string
	TypeFunctionHeader string
}{
	IncludeImport: true,

	SynopsisHeader: "###",
	HeadifySynopsis: true,

	ConstantHeader: "####",
	VariableHeader: "####",
	FunctionHeader: "####",
	TypeHeader: "####",
	TypeFunctionHeader: "####",
}

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
	if !Style.HeadifySynopsis {
		return target
	}
	return synopsisHeading_Regexp.ReplaceAllStringFunc(target, func(heading string) string {
		return fmt.Sprintf("%s %s", Style.SynopsisHeader, heading)
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
	return indent_Regexp.ReplaceAllString(target, indent + "$1")
}

func writeConstantSection(writer io.Writer, list []*doc.Value) {
	for _, entry := range list {
		fmt.Fprintf(writer, "%s\n%s\n", indentCode(sourceOfNode(entry.Decl)), formatIndent(entry.Doc))
	}
}

func writeVariableSection(writer io.Writer, list []*doc.Value) {
	for _, entry := range list {
		fmt.Fprintf(writer, "%s\n%s\n", indentCode(sourceOfNode(entry.Decl)), formatIndent(entry.Doc))
	}
}

func writeFunctionSection(writer io.Writer, list []*doc.Func, inTypeSection bool) {

	header := Style.FunctionHeader
	if inTypeSection {
		header = Style.TypeFunctionHeader
	}

	for _, entry := range list {
		receiver := " "
		if entry.Recv != "" {
			receiver = fmt.Sprintf("(%s) ", entry.Recv)
		}
		fmt.Fprintf(writer, "%s func %s%s\n\n%s\n%s\n", header, receiver, entry.Name, indentCode(sourceOfNode(entry.Decl)), formatIndent(entry.Doc))
	}
}

func writeTypeSection(writer io.Writer, list []*doc.Type) {

	header := Style.TypeHeader

	for _, entry := range list {
		fmt.Fprintf(writer, "%s type %s\n\n%s\n\n%s\n", header, entry.Name, indentCode(sourceOfNode(entry.Decl)), formatIndent(entry.Doc))
		writeConstantSection(writer, entry.Consts)
		writeVariableSection(writer, entry.Vars)
		writeFunctionSection(writer, entry.Funcs, true)
		writeFunctionSection(writer, entry.Methods, true)
	}
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

	dotImport := ""
	if read, err := ioutil.ReadFile(filepath.Join(path, ".import")); err == nil {
		dotImport = strings.TrimSpace(strings.Split(string(read), "\n")[0])
	}

	var buffer bytes.Buffer

	found := false
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
			// path is used to get the containing directory in the case of
			// command documentation
			path, err := filepath.Abs(path)
			if err != nil {
				panic(err)
			}
			_, name = filepath.Split(path)
			isCommand = true
		default:
			name = pkg.Name
			// Just a regular package
		}

		found = true
		document := &_document{
			name: name,
			pkg: pkg,
			isCommand: isCommand,
		}

		// Header
		fmt.Fprintf(&buffer, "# %s\n--\n", document.name)

		if !document.isCommand {
			// Import
			if Style.IncludeImport {
				if (dotImport != "") {
					fmt.Fprintf(&buffer, space(4) + "import \"%s\"\n\n", dotImport)
				}
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
			writeFunctionSection(&buffer, document.pkg.Funcs, false)

			// Type Section
			writeTypeSection(&buffer, document.pkg.Types)
		}

		break
	}

	if !found {
		rootPath, _ := filepath.Abs(path)
		fmt.Fprintf(os.Stderr, "No package/documentation found in %s (%s)\n", path, rootPath)
		os.Exit(64)
	}

	if debug {
		return
	}

	fmt.Println(strings.TrimSpace(buffer.String()))

	if *signature_flag {
		fmt.Printf("\n--\n**godocdown** http://github.com/robertkrimen/godocdown\n")
	}
}
