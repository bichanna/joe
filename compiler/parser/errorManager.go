package parser

import (
	"fmt"

	"github.com/bichanna/joe/compiler/tokenizer"
)

type ErrorType int

const (
	ErrUnexpectedSymbol ErrorType = iota
	ErrIllegalNumberFormat
	ErrUnexpectedEOF
	ErrExpectedStrLiteralEOF
	ErrIllegalStrFormat
	ErrExpectedCharLiteralEOF
	ErrIllegalCharFormat
	ErrGeneric
	ErrIllegalAccessDeclaration
	ErrIllegalBracketMisMatch
	ErrMissingBracket
	ErrInvalidAccessSpecifier
	ErrMultipleDefinition
	ErrPreviouslyDefined
	ErrDuplicateClass
	ErrRedundantToken
	ErrInternal
	ErrCouldNotResolve
	ErrExpectedRefOfType
	ErrInvalidCast
	ErrRedundantCast
	ErrRedundantImport
	ErrInvalidAccess
	ErrUnexpectedToken
	ErrSymbolAlreadyDefined
	ErrInvalidParam
	ErrIncompatibleTypes
	ErrDuplicateDeclaration

	NoErr ErrorType = 999
)

var predefinedErrs []KeyPair

// InitializeErrors initializes all errors.
func InitializeErrors() {
	predefinedErrs = append(
		predefinedErrs,
		*newKeyPair(ErrUnexpectedSymbol, "unexpected symbol"),
		*newKeyPair(ErrIllegalNumberFormat, "illegal number format"),
		*newKeyPair(ErrUnexpectedEOF, "unexpected end of file"),
		*newKeyPair(ErrExpectedStrLiteralEOF, "expected string literal before end of file"),
		*newKeyPair(ErrIllegalStrFormat, "illegal string format"),
		*newKeyPair(ErrExpectedCharLiteralEOF, "expected character literal before end of file"),
		*newKeyPair(ErrIllegalCharFormat, "illegal character literal format"),
		*newKeyPair(ErrGeneric, ""),
		*newKeyPair(ErrIllegalAccessDeclaration, "illegal specification of access specifier(s)"),
		*newKeyPair(ErrIllegalBracketMisMatch, "illegal symbol mismatch, unexpected bracket"),
		*newKeyPair(ErrMissingBracket, "missing bracket"),
		*newKeyPair(ErrInvalidAccessSpecifier, "invalid access specifier"),
		*newKeyPair(ErrPreviouslyDefined, ""),
		*newKeyPair(ErrDuplicateDeclaration, "duplicate class"),
		*newKeyPair(ErrRedundantToken, "redundant token"),
		*newKeyPair(ErrInternal, "internal runtime error (not your fault)"),
		*newKeyPair(ErrCouldNotResolve, "could not resolve symbol"),
		*newKeyPair(ErrExpectedRefOfType, "expected reference of type"),
		*newKeyPair(ErrInvalidCast, "invalid type cast"),
		*newKeyPair(ErrRedundantCast, "redundant type cast"),
		*newKeyPair(ErrRedundantImport, "redundant self import of module"),
		*newKeyPair(ErrUnexpectedToken, "unexpected token"),
		*newKeyPair(ErrInvalidAccess, "invalid access of"),
		*newKeyPair(ErrSymbolAlreadyDefined, ""),
		*newKeyPair(ErrInvalidParam, "invalid parameter of type"),
		*newKeyPair(ErrIncompatibleTypes, "incompatible types"),
		*newKeyPair(ErrDuplicateDeclaration, "duplicate declaration of"),
	)
}

type KeyPair struct {
	Key   ErrorType
	Value string
}

type ParserError struct {
	Id      ErrorType
	Error   string
	Line    uint
	Col     uint
	Warning bool
}

// ErrorManager per one file
type ErrorManager struct {
	fileName         string
	asIs             bool
	lines            []string
	aggresive        bool
	errors           []*ParserError
	warnings         []*ParserError
	unfilteredErrors []*ParserError
	possibleErrors   [][]*ParserError
	lastError        *ParserError
	lastCheckedError *ParserError
	teCurser         int64
	haveFoundErr     bool
	cm               bool
}

// NewErrorManager creates a new ErrorManager with various flags.
func NewErrorManager(lines []string, fileName string, asIs bool, aggresive bool) *ErrorManager {
	return &ErrorManager{
		lines:        lines,
		fileName:     fileName,
		asIs:         asIs,
		aggresive:    aggresive,
		teCurser:     0,
		haveFoundErr: false,
		cm:           false,
	}
}

// PrintErrors prints errors.
func (m *ErrorManager) PrintErrors() {
	if !m.asIs {
		if m.haveFoundErr {
			if m.aggresive {
				fmt.Print(m.getErrors(&m.unfilteredErrors))
			} else {
				fmt.Print(m.getErrors(&m.errors))
			}
		}
		fmt.Print(m.getErrors(&m.warnings))
	}
}

// GetErrorCount returns the number of solid errors (not unfiltered).
func (m *ErrorManager) GetErrorCount() uint64 {
	return 0
}

// GetWarningCount returns the number of warnings.
func (m *ErrorManager) GetWarningCount() uint64 {
	return 0
}

// GetUnfilteredErrorCount returns the number of all errors.
func (m *ErrorManager) GetUnfilteredErrorCount() uint64 {
	return 0
}

// CreateNewErrorWithToken creates a new error with a token.
func (m *ErrorManager) CreateNewErrorWithToken(err ErrorType, token *tokenizer.TokenEntity, extraComments string) int {
	return 0
}

// CreateNewErrorWithAST creates a new error with an AST.
func (m *ErrorManager) CreateNewErrorWithAST(err ErrorType, ast *AST, extraComments string) int {
	return 0
}

// CreateNewErrorWithLine creates a new error with the given line and column numbers.
func (m *ErrorManager) CreateNewErrorWithLineAndCol(err ErrorType, line, col, uint, extraComments string) {

}

// CreateNewWarningWithLineAndCol creates a new warning with the given line and column numbers.
func (m *ErrorManager) CreateNewWarningWithLineAndCol(err ErrorType, line, col, uint, extraComments string) {

}

// CreateNewWarningWithAST creates a new warning with the given AST.
func (m *ErrorManager) CreateNewWarningWithAST(err ErrorType, ast *AST, extraComments string) {

}

// HasErrors checks whether the manager has any errors or not.
func (m *ErrorManager) HasErrors() bool {
	return false
}

// HasWarnings checks whether the manager has any warnings or not.
func (m *ErrorManager) HasWarnings() bool {
	return false
}

// EnableErrorCheckMode enables error check mode.
func (m *ErrorManager) EnableErrorCheckMode() {

}

func (m *ErrorManager) Pass() {

}

func (m *ErrorManager) Fail() {

}

// GetLine returns the corresponding line.
func (m *ErrorManager) GetLine(line uint) string {
	return ""
}

// getErrorByID gets the specified error type.
func (m *ErrorManager) getErrorByID(err ErrorType) *KeyPair {
	return nil
}

func (m *ErrorManager) getPossibleErrorList() []*ParserError {
	return nil
}

func (m *ErrorManager) addPossibleError() {

}

func (m *ErrorManager) removePossibleError() {

}

// shouldReport checks if the given error should be reporeted or not.
func (m *ErrorManager) shouldReport(token *tokenizer.TokenEntity, lastErr *ParserError, e *ParserError) bool {
	return false
}

// getErrors combines errors and return error message as one giant error message.
func (m *ErrorManager) getErrors(errs *[]*ParserError) string {
	errMsg := ""
	for _, err := range *errs {
		if err.Warning {
			errMsg += fmt.Sprintf("%s:%d:%d: warning E20%d: %s\n", m.fileName, err.Line, err.Col, err.Id, err.Error)
		} else {
			errMsg += fmt.Sprintf("%s:%d:%d: error $50%d: %s\n", m.fileName, err.Line, err.Col, err.Id, err.Error)
		}
		errMsg += fmt.Sprintf("\t%s\n\t", m.GetLine(err.Line))

		for i := 0; i < int(err.Col)-1; i++ {
			errMsg += " "
		}
		errMsg += "^"
	}
	return errMsg
}

// printError directly prints a specific error.
func (m *ErrorManager) printError(err *ParserError) {
	if err.Warning {
		fmt.Printf("%s:%d:%d: warning E20%d: %s\n", m.fileName, err.Line, err.Col, err.Id, err.Error)
	} else {
		fmt.Printf("%s:%d:%d: error $50%d: %s\n", m.fileName, err.Line, err.Col, err.Id, err.Error)
	}
	fmt.Printf("\t%s\n\t", m.GetLine(err.Line))

	for i := 0; i < int(err.Col)-1; i++ {
		fmt.Print(" ")
	}
	fmt.Println("^")
}

// Check if the error is in the given error list.
func (m *ErrorManager) hasError(errs *[]*ParserError, e *ParserError) bool {
	return false
}

func (m *ErrorManager) shouldReportWarning(token *tokenizer.TokenEntity, lastErr *ParserError, err *ParserError) bool {
	return false
}

// newKeyPair creates a new KeyPair.
func newKeyPair(key ErrorType, value string) *KeyPair {
	return &KeyPair{
		Key:   key,
		Value: value,
	}
}

// newParserErrorLineAndCol creates a new ParserError with the given line and column numbers.
func newParserErrorLineAndCol(err KeyPair, l, c uint, addon string, warning bool) *ParserError {
	msg := err.Value
	if addon != "" {
		msg += ": " + addon
	}
	return &ParserError{
		Id:      err.Key,
		Error:   msg,
		Line:    l,
		Col:     c,
		Warning: warning,
	}
}

// newParserErrorFromToken creates a new ParserError with the line and column numbers taken from the given token.
func newParserErrorFromToken(err KeyPair, token *tokenizer.TokenEntity, addon string) *ParserError {
	msg := err.Value
	if addon != "" {
		msg += ": " + addon
	}
	return &ParserError{
		Id:      err.Key,
		Error:   msg,
		Line:    token.Line,
		Col:     token.Col,
		Warning: false,
	}
}
