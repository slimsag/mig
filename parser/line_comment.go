package parser

import (
	"io"
	"strings"

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
		return nil, p.error(ExpectedSingleLineComment)
	} else if ioErr != nil {
		return nil, p.ioError(ioErr)
	}
	if string(b) != "//" {
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
		b, ioErr := p.readByte()
		if ioErr == io.EOF {
			break
		} else if ioErr != nil {
			return nil, p.ioError(ioErr)
		}
		body = append(body, b)

		// Attempt parsing a newline.
		newLine, err := p.parseNewline()
		if err != nil {
			return nil, err
		}
		if newLine != nil {
			body = append(body, []byte(newLine.Body)...)
			break
		}
	}
	c.Body = string(body)

	// Parse post-body whitespace.
	postSpace := len(strings.TrimRight(c.Body, " \t\r\n"))
	if len(c.Body)-postSpace > 0 {
		c.PostBody = &ast.Whitespace{Body: c.Body[postSpace:]}
		c.Body = c.Body[:postSpace]
	}
	return c, nil
}
