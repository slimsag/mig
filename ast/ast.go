package ast

import "fmt"

// Module represents a module of code (a single directory of code).
type Module struct {
	// Name is the name of the module (the directory name).
	Name string

	// Files of code in the directory.
	Files []*File
}

// File represents a single file of code.
type File struct {
	// Name is the filename / path.
	Name string

	// Children nodes of the file. i.e. the immediate code found in the file.
	Children []Node
}

// SingleLineComment represents a single line comment.
type SingleLineComment struct {
	// PreOpen is indention before the opening `//`.
	PreOpen *Whitespace

	// PostOpen is indention after the opening `//`.
	PostOpen *Whitespace

	// Body is the rest of the line.
	Body string

	// PostBody is whitespace after the `// comment body`, including the
	// newline.
	PostBody *Whitespace
}

func (s *SingleLineComment) String() string {
	return fmt.Sprintf("SingleLineComment{PreOpen:%+v PostOpen:%+v Body:%+v PostBody: %+v}", s.PreOpen, s.PostOpen, s.Body, s.PostBody)
}

// Whitespace is one or more spaces, tabs, newlines or carriage returns in any
// sequence.
type Whitespace struct {
	Body string
}

func (w *Whitespace) String() string {
	return fmt.Sprintf("%q", w.Body)
}

// Node is any of the following:
//
//  *SingleLineComment
//
type Node interface{}
