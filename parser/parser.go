package parser

import (
	"io"

	"github.com/slimsag/mig/ast"
)

const (
	ZeroErrorCode ErrorCode = iota
	ExpectedFileBody
	ExpectedSingleLineComment
	ExpectedMultiLineComment
)

type ErrorCode int

func (e ErrorCode) String() string {
	switch e {
	case ExpectedFileBody:
		return "expected file body"
	case ExpectedSingleLineComment:
		return "expected single line comment"
	case ExpectedMultiLineComment:
		return "expected multi-line comment"
	}
	panic("invalid error code")
}

// Error represents a parser error.
type Error struct {
	// Start and End bytes of where the error is located.
	Start, End int

	// IO error, if any.
	IOError error

	// Error code.
	Code ErrorCode
}

func (e Error) Error() string {
	if e.IOError != nil {
		return e.IOError.Error()
	}
	return e.Code.String()
}

// Parse parses the code from buf, as if it was part of the given module and in
// the named file.
func Parse(mod *ast.Module, filename string, buf []byte) error {
	file, err := (&parser{
		mod: mod,
		buf: buf,
	}).parseFile(filename, buf)
	if err != nil {
		return err
	}
	mod.Files = append(mod.Files, file)
	return nil
}

type parser struct {
	mod                *ast.Module
	readHead           int
	buf                []byte
	current            []ast.Node
	foundNonWhitespace bool
}

func (p *parser) error(c ErrorCode) *Error {
	return &Error{
		Code: c,
	}
}

func (p *parser) ioError(err error) *Error {
	return &Error{IOError: err}
}

func (p *parser) unreadBytes(n int) {
	p.readHead -= n
	if p.readHead < 0 {
		panic("cannot unread to < 0")
	}
}

func (p *parser) peek(n int) ([]byte, error) {
	if len(p.buf) < p.readHead+n {
		return nil, io.EOF
	}
	return p.buf[p.readHead : p.readHead+n], nil
}

func (p *parser) next(n int) ([]byte, error) {
	if len(p.buf) < p.readHead+n {
		return nil, io.EOF
	}
	data := p.buf[p.readHead : p.readHead+n]
	p.readHead += n
	return data, nil
}

func (p *parser) readByte() (byte, error) {
	data, err := p.next(1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (p *parser) emit(n ast.Node) {
	if _, isWhitespace := n.(*ast.Whitespace); !isWhitespace {
		p.foundNonWhitespace = true
	}
	p.current = append(p.current, n)
}

func (p *parser) parseFile(name string, body []byte) (*ast.File, error) {
	file := &ast.File{
		Name: name,
	}
	for {
		// Module bodies may contain single line comments.
		singleLineComment, err := p.parseSingleLineComment()
		if err != nil && err.Code != ExpectedSingleLineComment {
			return nil, err
		} else if err == nil {
			p.emit(singleLineComment)
			continue
		}

		// Module bodies may contain multi-line comments.
		multiLineComment, err := p.parseMultiLineComment()
		if err != nil && err.Code != ExpectedMultiLineComment {
			return nil, err
		} else if err == nil {
			p.emit(multiLineComment)
			continue
		}

		break
	}

	// Check that we found something other than just whitespace.
	if !p.foundNonWhitespace {
		return nil, p.error(ExpectedFileBody)
	}
	file.Children = p.current
	p.current = nil
	return file, nil
}

func (p *parser) parseSpacesOrTabs() (*ast.Whitespace, *Error) {
	any, err := p.parseAny(' ', '\t')
	if err != nil {
		return nil, err
	}
	if len(any) == 0 {
		return nil, nil
	}
	return &ast.Whitespace{Body: string(any)}, nil
}

func (p *parser) parseNewline() (*ast.Whitespace, *Error) {
	any, err := p.parseAny('\r', '\n')
	if err != nil {
		return nil, err
	}
	if len(any) == 0 {
		return nil, nil
	}
	return &ast.Whitespace{Body: string(any)}, nil
}

func (p *parser) parseSpaceAndNewline() (*ast.Whitespace, *Error) {
	space, err := p.parseSpacesOrTabs()
	if err != nil {
		return nil, err
	}
	newLine, err := p.parseNewline()
	if err != nil {
		return nil, err
	}
	if space == nil && newLine == nil {
		return nil, nil
	}
	w := &ast.Whitespace{}
	if space != nil {
		w.Body += space.Body
	}
	if newLine != nil {
		w.Body += newLine.Body
	}
	return w, nil
}

// parseAny tries to parse any of the characters in set, aborting if any
// character not found in set is encountered.
func (p *parser) parseAny(set ...byte) ([]byte, *Error) {
	var body []byte
	for {
		b, ioErr := p.readByte()
		if ioErr == io.EOF {
			break
		} else if ioErr != nil {
			return nil, p.ioError(ioErr)
		}
		shouldBreak := true
		for _, c := range set {
			if b == c {
				shouldBreak = false
				break
			}
		}
		if shouldBreak {
			p.unreadBytes(1)
			break
		}
		body = append(body, b)
	}
	return body, nil
}
