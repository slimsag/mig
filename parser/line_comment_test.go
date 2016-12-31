package parser_test

import (
	"reflect"
	"testing"

	"github.com/slimsag/mig/ast"
	"github.com/slimsag/mig/parser"
)

// TestParseLineComments tests that parsing line comments works.
func TestParseLineComments(t *testing.T) {
	tests := []struct {
		Name string
		Code string
		Want []ast.Node
	}{
		{
			Name: "empty",
			Code: "//",
			Want: []ast.Node{&ast.SingleLineComment{}},
		},
		{
			Name: "no_space",
			Code: "//foobar",
			Want: []ast.Node{&ast.SingleLineComment{Body: "foobar"}},
		},
		{
			Name: "std_space",
			Code: "// foobar",
			Want: []ast.Node{&ast.SingleLineComment{
				PostOpen: &ast.Whitespace{Body: " "},
				Body:     "foobar",
			}},
		},
		{
			Name: "prefix_whitespace",
			Code: "// \t foobar",
			Want: []ast.Node{&ast.SingleLineComment{
				PostOpen: &ast.Whitespace{Body: " \t "},
				Body:     "foobar",
			}},
		},
		{
			Name: "suffix_whitespace_lf",
			Code: "// foobar \n",
			Want: []ast.Node{&ast.SingleLineComment{
				PostOpen: &ast.Whitespace{Body: " "},
				Body:     "foobar",
				PostBody: &ast.Whitespace{Body: " \n"},
			}},
		},
		{
			Name: "suffix_whitespace_cr",
			Code: "// foobar \r",
			Want: []ast.Node{&ast.SingleLineComment{
				PostOpen: &ast.Whitespace{Body: " "},
				Body:     "foobar",
				PostBody: &ast.Whitespace{Body: " \r"},
			}},
		},
		{
			Name: "suffix_whitespace_crlf",
			Code: "// foobar \r\n",
			Want: []ast.Node{&ast.SingleLineComment{
				PostOpen: &ast.Whitespace{Body: " "},
				Body:     "foobar",
				PostBody: &ast.Whitespace{Body: " \r\n"},
			}},
		},
		{
			Name: "multiple_lf",
			Code: "// foobar \n// baz\n",
			Want: []ast.Node{
				&ast.SingleLineComment{
					PostOpen: &ast.Whitespace{Body: " "},
					Body:     "foobar",
					PostBody: &ast.Whitespace{Body: " \n"},
				},
				&ast.SingleLineComment{
					PostOpen: &ast.Whitespace{Body: " "},
					Body:     "baz",
					PostBody: &ast.Whitespace{Body: "\n"},
				},
			},
		},
		{
			Name: "multiple_cr",
			Code: "// foobar \r// baz\r",
			Want: []ast.Node{
				&ast.SingleLineComment{
					PostOpen: &ast.Whitespace{Body: " "},
					Body:     "foobar",
					PostBody: &ast.Whitespace{Body: " \r"},
				},
				&ast.SingleLineComment{
					PostOpen: &ast.Whitespace{Body: " "},
					Body:     "baz",
					PostBody: &ast.Whitespace{Body: "\r"},
				},
			},
		},
		{
			Name: "multiple_crlf",
			Code: "// foobar \r\n// baz\r\n",
			Want: []ast.Node{
				&ast.SingleLineComment{
					PostOpen: &ast.Whitespace{Body: " "},
					Body:     "foobar",
					PostBody: &ast.Whitespace{Body: " \r\n"},
				},
				&ast.SingleLineComment{
					PostOpen: &ast.Whitespace{Body: " "},
					Body:     "baz",
					PostBody: &ast.Whitespace{Body: "\r\n"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mod := &ast.Module{Name: "dirname"}
			err := parser.Parse(mod, "filename.mg", []byte(test.Code))
			if err != nil {
				t.Fatal(err)
			}
			if len(mod.Files[0].Children) != len(test.Want) {
				t.Fatalf("got %v children, expected %v\n", len(mod.Files[0].Children), len(test.Want))
			}
			for i, got := range mod.Files[0].Children {
				want := test.Want[i]
				if !reflect.DeepEqual(got, want) {
					t.Logf("got  %+v\n", got)
					t.Fatalf("want %+v\n", want)
				}
			}
		})
	}
}
