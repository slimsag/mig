package parser

import (
	"io"

	"github.com/slimsag/mig/ast"
)

// parseSingleLineComment first parses preceding whitespace, then a comment
// string.
func (p *parser) parseSingleLineComment() (c *ast.SingleLineComment, err *Error) {
	c = &ast.SingleLineComment{}

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
		return nil, p.error(ExpectedSingleLineComment)
	} else if ioErr != nil {
		return nil, p.ioError(ioErr)
	}
	if string(b) != "//" {
		if c.PreOpen != nil {
			p.unreadBytes(len(c.PreOpen.Body))
		}
		p.unreadBytes(2)
		return nil, p.error(ExpectedSingleLineComment)
	}

	// Parse post-open whitespace.
	c.PostOpen, err = p.parseSpacesOrTabs()
	if err != nil {
		return nil, err
	}

	// Consume comment body.
	var body []byte
	for {
		// Attempt parsing whitespace + a newline.
		c.PostBody, err = p.parseSpaceAndNewline()
		if err != nil {
			return nil, err
		}
		if c.PostBody != nil {
			break
		}

		b, ioErr := p.readByte()
		if ioErr == io.EOF {
			break
		} else if ioErr != nil {
			return nil, p.ioError(ioErr)
		}
		body = append(body, b)
	}
	c.Body = string(body)
	return c, nil
}
