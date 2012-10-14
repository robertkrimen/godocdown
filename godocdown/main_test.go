package main

import (
	"testing"
	. "github.com/robertkrimen/terst"
)

func TestIndent(t *testing.T) {
	Terst(t)

	Is(indent("1\n  2\n\n  3\n  4\n", "  "), "  1\n    2\n\n    3\n    4\n")
}
