package parser

import (
	"io"

	"github.com/slimsag/mig/ast"
)

// parseMultiLineComment parses a multi-line comment.
func (p *parser) parseMultiLineComment() (c *ast.MultiLineComment, err *Error) {
	c = &ast.MultiLineComment{}

	// Parse pre-open spaces or tabs.
	c.PreOpen, err = p.parseSpacesOrTabs()
	if err != nil {
		p.unreadBytes(len(c.PreOpen.Body))
		return nil, err
	}

	// Peek for a comment string.
	b, ioErr := p.next(2)
	if ioErr == io.EOF {
		if c.PreOpen != nil {
			p.unreadBytes(len(c.PreOpen.Body))
		}
		return nil, p.error(ExpectedMultiLineComment)
	} else if ioErr != nil {
		return nil, p.ioError(ioErr)
	}
	if string(b) != "/*" {
		if c.PreOpen != nil {
			p.unreadBytes(len(c.PreOpen.Body))
		}
		p.unreadBytes(2)
		return nil, p.error(ExpectedMultiLineComment)
	}

	// Consume comment body.
	var body []byte
	for {
		b, ioErr := p.readByte()
		if ioErr == io.EOF {
			break
		} else if ioErr != nil {
			return nil, p.ioError(ioErr)
		}
		if b == '*' {
			b2, ioErr := p.next(1)
			if ioErr == io.EOF {
				return nil, p.error(ExpectedMultiLineComment)
			} else if ioErr != nil {
				return nil, p.ioError(ioErr)
			}
			if b2[0] == '/' {
				break
			}
		}
		body = append(body, b)
	}
	c.Body = string(body)

	// Parse post-body whitespace.
	c.PostBody, err = p.parseSpaceAndNewline()
	if err != nil {
		return nil, err
	}
	return c, nil
}
