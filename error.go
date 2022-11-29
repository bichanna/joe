package main

import (
	"fmt"
	"strings"
)

type ErrorManager struct {
	filename   string
	unfiltered []*ParseError
	filtered   []*ParseError
	aggressive bool
}

// newErrorManager creates a new Error Manager.
func newErrorManager(filename string, aggressive bool) *ErrorManager {
	return &ErrorManager{
		filename:   filename,
		aggressive: aggressive,
	}
}

type ParseError struct {
	msg  string
	line int
	col  int
}

// newError creates a new ParseError with given line and column numbers.
func (m *ErrorManager) newError(line, col int, msg string) {
	err := &ParseError{
		msg: msg,
		line: line,
		col: col,
	}

	m.unfiltered = append(m.unfiltered, err)
}

// newErrorWithPosition creates a new ParseError with the given Position.
func (m *ErrorManager) newErrorWithPosition(pos *Position, msg string) {
	m.newError(pos.line, pos.col, msg)
}

// newErrorWithToken creates a new ParseError from a token.
func (m *ErrorManager) newErrorWithToken(token *Token, msg string) {
	m.newError(token.pos.line, token.pos.col, msg)
}

// formatError formats a ParseError.
func (m *ErrorManager) formatError(e *ParseError) string {
	return fmt.Sprintf("\n%s:%d:%d: SyntaxError: %s", m.filename, e.line, e.col, e.msg)
}

// stringifyErrors stringifies captured errors.
func (m *ErrorManager) stringifyErrors() string {
	var (
		errors  []*ParseError
		builder = strings.Builder{}
	)
	builder.WriteString("aggressive error reporting: ")

	// if aggressive error reporting is enabled, show all errors
	if m.aggressive {
		errors = m.unfiltered
		builder.WriteString("disabled\n")
	} else {
		errors = m.filtered
		builder.WriteString("enabled\n")
	}

	// concatenate all the errors
	for _, err := range errors {
		builder.WriteString(m.formatError(err))
	}

	return builder.String()
}

// reportErrors reports errors.
func (m *ErrorManager) reportErrors() {
	if m.aggressive {
		// TODO: filter the unfiltered errors
		m.filtered = append(m.filtered, m.unfiltered...)
	}

	fmt.Println(m.stringifyErrors())
}
