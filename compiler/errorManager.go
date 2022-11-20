package compiler

import "fmt"

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
		lines:            lines,
		fileName:         fileName,
		asIs:             asIs,
		aggresive:        aggresive,
		teCurser:         0,
		haveFoundErr:     false,
		cm:               false,
		lastError:        &ParserError{},
		lastCheckedError: &ParserError{},
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
	return uint64(len(m.errors))
}

// GetWarningCount returns the number of warnings.
func (m *ErrorManager) GetWarningCount() uint64 {
	return uint64(len(m.warnings))
}

// GetUnfilteredErrorCount returns the number of all errors.
func (m *ErrorManager) GetUnfilteredErrorCount() uint64 {
	return uint64(len(m.unfilteredErrors))
}

// CreateNewErrorWithToken creates a new error with a token.
func (m *ErrorManager) CreateNewErrorWithToken(err ErrorType, token *TokenEntity, extraComments string) int {
	kp := m.getErrorByID(err)
	newErr := newParserErrorFromToken(kp, token, extraComments)
	var lastError *ParserError
	if m.cm {
		lastError = m.lastCheckedError
	} else {
		lastError = m.lastError
	}

	if m.shouldReport(nil, lastError, newErr) || (m.aggresive && m.asIs) {
		if m.asIs {
			m.printError(newErr)
		} else if m.cm {
			*m.getPossibleErrorList() = append(*m.getPossibleErrorList(), newErr)
			m.lastCheckedError = newErr
			return 1
		}
		m.haveFoundErr = true
		m.errors = append(m.errors, newErr)
		m.unfilteredErrors = append(m.unfilteredErrors, newErr)
		m.lastError = newErr
		return 1
	} else {
		m.unfilteredErrors = append(m.unfilteredErrors, newErr)
	}
	return 0
}

// CreateNewErrorWithAST creates a new error with an AST.
func (m *ErrorManager) CreateNewErrorWithAST(err ErrorType, ast *AST, extraComments string) int {
	kp := m.getErrorByID(err)
	newErr := newParserErrorLineAndCol(kp, ast.Line, ast.Col, extraComments, false)
	var lastError *ParserError
	if m.cm {
		lastError = m.lastCheckedError
	} else {
		lastError = m.lastError
	}

	if m.shouldReport(nil, lastError, newErr) || (m.aggresive && m.asIs) {
		if m.asIs {
			m.printError(newErr)
		} else if m.cm {
			*m.getPossibleErrorList() = append(*m.getPossibleErrorList(), newErr)
			m.lastCheckedError = newErr
			return 1
		}
		m.haveFoundErr = true
		m.errors = append(m.errors, newErr)
		m.unfilteredErrors = append(m.unfilteredErrors, newErr)
		m.lastError = newErr
		return 1
	} else {
		m.unfilteredErrors = append(m.unfilteredErrors, newErr)
	}

	return 0
}

// CreateNewErrorWithLine creates a new error with the given line and column numbers.
func (m *ErrorManager) CreateNewErrorWithLineAndCol(err ErrorType, line, col uint, extraComments string) {
	kp := m.getErrorByID(err)
	newErr := newParserErrorLineAndCol(kp, line, col, extraComments, false)
	var lastError *ParserError
	if m.cm {
		lastError = m.lastCheckedError
	} else {
		lastError = m.lastError
	}

	if m.shouldReport(nil, lastError, newErr) || (m.aggresive && m.asIs) {
		if m.asIs {
			m.printError(newErr)
		} else if m.cm {
			*m.getPossibleErrorList() = append(*m.getPossibleErrorList(), newErr)
			m.lastCheckedError = newErr
			return
		}
		m.haveFoundErr = true
		m.errors = append(m.errors, newErr)
		m.unfilteredErrors = append(m.unfilteredErrors, newErr)
		m.lastError = newErr
	} else {
		m.unfilteredErrors = append(m.unfilteredErrors, newErr)
	}
}

// CreateNewWarningWithLineAndCol creates a new warning with the given line and column numbers.
func (m *ErrorManager) CreateNewWarningWithLineAndCol(err ErrorType, line, col uint, extraComments string) {
	kp := m.getErrorByID(err)
	newWrn := newParserErrorLineAndCol(kp, line, col, extraComments, true)
	var lastWrn *ParserError
	if len(m.warnings) > 0 {
		lastWrn = m.warnings[len(m.warnings)-1]
	} else {
		if m.cm {
			lastWrn = m.lastCheckedError
		} else {
			lastWrn = m.lastError
		}
	}

	if len(m.warnings) == 0 || m.shouldReportWarning(nil, lastWrn, newWrn) {
		if m.asIs {
			m.printError(newWrn)
		}
		m.warnings = append(m.warnings, newWrn)
	}
}

// CreateNewWarningWithAST creates a new warning with the given AST.
func (m *ErrorManager) CreateNewWarningWithAST(err ErrorType, ast *AST, extraComments string) {
	m.CreateNewWarningWithLineAndCol(err, ast.Line, ast.Col, extraComments)
}

// hasErrors checks whether the manager has any errors or not.
func (m *ErrorManager) hasErrors() bool {
	return m.haveFoundErr && len(m.unfilteredErrors) != 0
}

// hasWarnings checks whether the manager has any warnings or not.
func (m *ErrorManager) hasWarnings() bool {
	return false
}

// enableErrorCheckMode enables error check mode.
func (m *ErrorManager) enableErrorCheckMode() {
	m.cm = true
	m.addPossbleErrorList()
}

func (m *ErrorManager) pass() {
	m.lastCheckedError = nil
	m.removePossibleErrorList()
}

func (m *ErrorManager) fail() {
	if len(m.possibleErrors) > 0 {
		for _, err := range *m.getPossibleErrorList() {
			if m.shouldReport(nil, m.lastError, err) {
				m.errors = append(m.errors, err)
				m.lastError = err
				m.unfilteredErrors = append(m.unfilteredErrors, err)
			}
		}

		if m.teCurser <= 0 {
			m.lastError = m.lastCheckedError
			m.haveFoundErr = true
		}
	}

	m.lastCheckedError = nil
	m.removePossibleErrorList()
}

// getLine returns the corresponding line.
func (m *ErrorManager) getLine(line uint) string {
	if line-1 >= uint(len(m.lines)) {
		return "EOF"
	} else {
		return m.lines[line-1]
	}
}

// getErrorByID gets the specified error type.
func (m *ErrorManager) getErrorByID(err ErrorType) *KeyPair {
	for _, i := range predefinedErrs {
		if i.Key == err {
			return &i
		}
	}
	return nil
}

// getPossibleErrorList returns the last element of possibleErrors.
func (m *ErrorManager) getPossibleErrorList() *[]*ParserError {
	return &m.possibleErrors[len(m.possibleErrors)-1]
}

// addPossbleErrorList adds a new element (a list of ParserErrors) to possibleErrors.
func (m *ErrorManager) addPossbleErrorList() {
	m.possibleErrors = append(m.possibleErrors, []*ParserError{})
	m.teCurser++
}

// removePossibleErrorList removes the last element of possibleErrors.
func (m *ErrorManager) removePossibleErrorList() {
	if len(m.possibleErrors) != 0 {
		m.possibleErrors = m.possibleErrors[:len(m.possibleErrors)-1]
		m.teCurser--
		if m.teCurser < 0 {
			m.cm = false
		}
	}
}

// shouldReport checks if the given error should be reported or not.
func (m *ErrorManager) shouldReport(token *TokenEntity, lastError *ParserError, e *ParserError) bool {
	if lastError.Error != e.Error && !m.hasError(&m.errors, e) && !(lastError.Line == e.Line && lastError.Col == e.Col) {
		if token != nil &&
			!(token.IsSingle() ||
				token.GetID() == CharLiteral ||
				token.GetID() == StringLiteral ||
				token.GetID() == IntegerLiteral) {
			return lastError.Line-e.Line != 1
		}
		return true
	}
	return false
}

// shouldReportWarning checks if the given warning should be reported or not. (exactly the same except for the name)
func (m *ErrorManager) shouldReportWarning(token *TokenEntity, lastErr *ParserError, err *ParserError) bool {
	return m.shouldReport(token, lastErr, err)
}

// getErrors combines errors and return error message as one giant error message.
func (m *ErrorManager) getErrors(errs *[]*ParserError) string {
	errMsg := ""
	for _, err := range *errs {
		if err.Warning {
			errMsg += fmt.Sprintf("%s:%d:%d: warning E20%d: %s\n", m.fileName, err.Line, err.Col, err.Id, err.Error)
		} else {
			errMsg += fmt.Sprintf("%s:%d:%d: error E50%d: %s\n", m.fileName, err.Line, err.Col, err.Id, err.Error)
		}
		errMsg += fmt.Sprintf("\t%s\n\t", m.getLine(err.Line))

		for i := 0; i < int(err.Col)-1; i++ {
			errMsg += " "
		}
		errMsg += "^\n"
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
	fmt.Printf("\t%s\n\t", m.getLine(err.Line))

	for i := 0; i < int(err.Col)-1; i++ {
		fmt.Print(" ")
	}
	fmt.Println("^")
}

// Check if the error is in the given error list.
func (m *ErrorManager) hasError(errs *[]*ParserError, e *ParserError) bool {
	for _, err := range *errs {
		if err.Error == e.Error {
			return true
		}
	}
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
func newParserErrorLineAndCol(err *KeyPair, l, c uint, addon string, warning bool) *ParserError {
	msg := err.Value + addon
	return &ParserError{
		Id:      err.Key,
		Error:   msg,
		Line:    l,
		Col:     c,
		Warning: warning,
	}
}

// newParserErrorFromToken creates a new ParserError with the line and column numbers taken from the given token.
func newParserErrorFromToken(err *KeyPair, token *TokenEntity, addon string) *ParserError {
	msg := err.Value + addon
	return &ParserError{
		Id:      err.Key,
		Error:   msg,
		Line:    token.Line,
		Col:     token.Col,
		Warning: false,
	}
}
