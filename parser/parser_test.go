package parser_test

import (
	"testing"

	"github.com/slimsag/mig/ast"
	"github.com/slimsag/mig/parser"
)

// TestParseEmpty tests that parsing an empty file returns ExpectedModuleBody.
func TestParseEmpty(t *testing.T) {
	err := parser.Parse(&ast.Module{Name: "dirname"}, "filename.mg", nil)
	if err.(*parser.Error).Code != parser.ExpectedFileBody {
		t.Fail()
	}
}
