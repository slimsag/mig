# Mig Programming Language Specification

- [x] [Code Files](#code-files)
- [ ] [Comments](#comments)
  - [x] [Single-Line Comments](#single-line-comments)
  - [ ] [Multi-Line Comments](#multi-line-comments)
- [ ] [Constants](#constants)
- [ ] [Variables](#variables)
- [ ] [Functions](#functions)
- [ ] [Line Endings](#line-endings)

## Code Files

Mig code files use the extension `.mg`, and the file path under `src` determines
the namespace of the code. For example:

```
mycode/src/people/john.mg
mycode/src/people/kim.mg
mycode/src/people/important/king.mg
```

1. In the above example, the file paths under `mycode/src` are considered the
namespace of the code.
  - A variable `Foobar` defined in `john.mg` would be accessible outside the
  file as `people.Foobar`.
  - A variable `Foobar` defined in `king.mg` would be accessible outside the
  file as `people.important.Foobar`.
2. Both `people/john.mg` and `people/kim.mg` share the same namespace.
3. `people/john.mg` and `people/important/king.mg` do not share the same
  namespace, because they are in separate directories.

## Comments

Mig supports two kinds of code comments, single-line comments and multi-line
comments.

Comments are in the [CommonMark](http://spec.commonmark.org) format and
thus can be automatically turned into rich HTML documentation for viewing in
the web.

### Single-Line Comments

A single-line comment can be written in the code using `//`. For example:

```
// # This is a comment!
```

The comment begins when `//` is encountered, and ends when the [line ends](#line-endings).

### Multi-Line Comments

Multi-line comments work exactly like line comments, except they can be used to
write a comment that spans across multiple lines. For example:

```
/*
  I span
  multiple lines!
*/
```

Both the beginning `/*` and ending `*/` must be _on their own line_ with _no
other text except whitespace_. They may be nested, for example:

```
/*
  A multi-line comment
  /*
    Inside of a
  */
  multi line comment!
*/
```

### Constants

A constant is just like a variable, it has a [name](#Names) and
[value](#Values), but unlike a variable the value cannot change. The syntax is
identical to [Variables]() except replacing the keyword `var` with `const`.

### Variables

A variable consists of a [name](#Names) and [value](#Values), the value is said
to be _variable_ because it can be changed after it is declared. For example:

```
var{
  x = 10            // signed 64-bit integer
  y uint32 = 10     // unsigned 32-bit integer
  z = 10.3          // 64-bit floating point number
  w = "hello world" // a string of characters
}
```

- The pattern used is `var <name> <type> = <value>` where `<type>` is optional.
- A single variable should be defined as `var x = 10`, while multiple variables
  should be defined in a `var{...}` block as shown above.

### Functions

Functions take any fixed number of parameters, and return a single value. For
example:

```
// foobar returns no value and takes no parameters.
func foobar() {
  // ...
}

// foobar returns no value and takes a single integer parameter.
func foobar(x int64) {
  // ...
}

// int_to_string returns a string value and takes a single integer parameter.
func int_to_string(x int64) string {
  // ...
}

// multiply returns a integer value and takes two integer parameter.
func multiply(x, y int64) int64 {
  // ...
}

// split returns two strings and takes two parameters.
func split(a string, b int64) (firstHalf, secondHalf string) {
  // ...
}
```

# Line Endings

Mig considers line feeds (`\n`), carriage returns (`\r`), and CRLF (`\r\n`) to be the end of a line. This enables it to work with code written in all standard Unix, Windows, etc text editors.
