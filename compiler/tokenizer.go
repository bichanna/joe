package compiler

import "unicode"

type Tokenizer struct {
	EOFToken   *TokenEntity
	fileName   string
	entities   []*TokenEntity
	errManager *ErrorManager
	lines      []string
	tokens     string
	empty      string
	cursor     uint
	length     uint
	line       uint
	col        uint
}

// NewTokenizer creates a new Tokenizer.
func NewTokenizer(tokens string, fileName string) *Tokenizer {
	tokenizer := &Tokenizer{
		length:   uint(len(tokens)),
		line:     1,
		fileName: fileName,
	}

	if tokens != "" {
		tokenizer.tokens = tokens
	}

	tokenizer.parse()

	return tokenizer
}

// getEntityCount gets the length of the entity list.
func (t *Tokenizer) GetEntityCount() int {
	return len(t.entities)
}

// getEntites gets the list of token entities.
func (t *Tokenizer) GetEntites() *[]*TokenEntity {
	return &t.entities
}

// getLines gets the lines.
func (t *Tokenizer) getLines() *[]string {
	return &t.lines
}

func (t *TokenEntity) GetData() *string {
	return nil
}

// parse parses the source file.
func (t *Tokenizer) parse() {
	if t.length == 0 {
		return
	}

	// get all the lines.
	t.parseLines()

	t.errManager = NewErrorManager(t.lines, t.fileName, false, COptionsAggressiveErrorReporting)
	t.EOFToken = NewTokenEntity(uint(len(t.lines)), 0, Single, EOF, "EOF")

	for {
	start:
	invalidate_whitespace:
		// check for comments
		for !t.isEnd() && isWhitespace(t.current()) {
			if t.current() == '\n' {
				t.newLine()
			}
			t.advance()
		}
		// check if there's any comments
		// mode == 1 -> one-line comment
		// mode == 2 -> multi-line comment
		// mode == 0 -> no comment
		var mode uint = 0
		if t.isEnd() || t.peekEnd(1) {
			goto scan
		} else if !commentStart(t.current(), t.peek(1), &mode) {
			goto scan
		}
		t.col += 2
		t.cursor += 2

		for !t.isEnd() && !commentEnd(t.current(), t.peek(1), &mode) {
			if t.current() == '\n' {
				t.newLine()
			}
			t.advance()
		}

		if t.current() == '\n' {
			t.newLine()
		}

		if !t.isEnd() {
			t.cursor += mode
			goto invalidate_whitespace
		} else if t.isEnd() && !commentEnd(t.current(), t.peek(1), &mode) {
			t.errManager.CreateNewErrorWithLineAndCol(ErrGeneric, t.line, t.col, "unterminated block comment")
		}

	scan:
		if t.isEnd() {
			goto end
		} else if t.isSymbol(t.current()) {
			if !t.peekEnd(1) { // two symbols: /=, +=, --, etc.
				var ttype TokenType = None
				chars := []byte{t.current(), t.peek(1)}

				if chars[0] == '<' && chars[1] == '=' {
					ttype = LTE
				} else if chars[0] == '>' && chars[1] == '=' {
					ttype = GTE
				} else if chars[0] == '!' && chars[1] == '=' {
					ttype = NotEq
				} else if chars[0] == '=' && chars[1] == '=' {
					ttype = EQEQ
				} else if chars[0] == '<' && chars[1] == '<' {
					ttype = SHL
				} else if chars[0] == '>' && chars[1] == '>' {
					ttype = SHR
				} else if chars[0] == '&' && chars[1] == '&' {
					ttype = And
				} else if chars[0] == '|' && chars[1] == '|' {
					ttype = Or
				} else if chars[0] == '-' && chars[1] == '>' {
					ttype = PTR
				} else if chars[0] == '+' && chars[1] == '+' {
					ttype = Inc
				} else if chars[0] == '-' && chars[1] == '-' {
					ttype = Dec
				} else if chars[0] == '+' && chars[1] == '=' {
					ttype = PlusEq
				} else if chars[0] == '*' && chars[1] == '=' {
					ttype = MultEq
				} else if chars[0] == '-' && chars[1] == '=' {
					ttype = MinusEq
				} else if chars[0] == '/' && chars[1] == '=' {
					ttype = DivEq
				} else if chars[0] == '&' && chars[1] == '=' {
					ttype = AndEq
				} else if chars[0] == '|' && chars[1] == '=' {
					ttype = OrEq
				} else if chars[0] == '^' && chars[1] == '=' {
					ttype = XOrEq
				} else if chars[0] == '%' && chars[1] == '=' {
					ttype = ModEq
				}

				if ttype != None {
					t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, ttype, string(chars)))
					t.cursor += 2
					goto start
				}
			}

			if t.current() == '<' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, LessThan, "<"))
			} else if t.current() == '>' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, GreaterThan, ">"))
			} else if t.current() == ';' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, SemiColon, ";"))
			} else if t.current() == ':' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Colon, ":"))
			} else if t.current() == '+' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Plus, "+"))
			} else if t.current() == '-' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Minus, "-"))
			} else if t.current() == '*' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Mult, "*"))
			} else if t.current() == ',' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Comma, ","))
			} else if t.current() == '=' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Assign, "="))
			} else if t.current() == '$' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Dollar, "$"))
			} else if t.current() == '!' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Not, "!"))
			} else if t.current() == '/' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Div, "/"))
			} else if t.current() == '%' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Mod, "%"))
			} else if t.current() == '(' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, LeftParen, "("))
			} else if t.current() == ')' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, RightParen, ")"))
			} else if t.current() == '{' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, LeftCurly, "{"))
			} else if t.current() == '}' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, RightCurly, "}"))
			} else if t.current() == '.' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Dot, "."))
			} else if t.current() == '[' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, LeftBrace, "["))
			} else if t.current() == ']' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, RightBrace, "]"))
			} else if t.current() == '&' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, BitAnd, "&"))
			} else if t.current() == '|' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, BitOr, "|"))
			} else if t.current() == '^' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, XOr, "^"))
			} else if t.current() == '?' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Question, "?"))
			} else if t.current() == '#' {
				t.entities = append(t.entities, NewTokenEntity(t.line, t.col, Single, Hash, "#"))
			} else {
				t.errManager.CreateNewErrorWithLineAndCol(ErrUnexpectedSymbol, t.line, t.col, "'"+string(t.current())+"'")
			}

			t.advance()
			goto start
		} else if isLetter(t.current()) || t.current() == '_' { // blah, int, david, _273, etc.
			variable := ""
			hasLetter := false

			for !t.isEnd() && (isLetter(t.current()) || isNumber(t.current()) || t.current() == '_') {
				if isLetter(t.current()) {
					hasLetter = true
				}
				variable += string(t.current())
				t.advance()
			}

			if !hasLetter {
				t.errManager.CreateNewErrorWithLineAndCol(ErrGeneric, t.line, t.col, "expected at least 1 letter in identifier '"+variable+"'")
			} else {
				t.entities = append(t.entities, NewDefaultTokenEntity(t.line, t.col, Identifier, variable))
			}

			goto start
		} else if isNumber(t.current()) { // 23.1234, 134_231_2444, etc.
			if t.current() == '0' && t.peek(1) == 'x' { // checking for a hex number
				number := "0x"
				underscoreOk := false

				t.col += 2
				t.cursor += 2
				for !t.isEnd() && isHexNum(t.current()) || t.current() == '_' {
					if isHexNum(t.current()) {
						underscoreOk = true
					} else {
						if !underscoreOk {
							t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalNumberFormat, t.line, t.col, ", unexpected or illegally placed underscore")
							break
						}

						t.advance()
						continue
					}

					number += string(t.current())
					t.advance()
				}

				t.entities = append(t.entities, NewDefaultTokenEntity(t.line, t.col, HexLiteral, number))
				goto start
			} else {
				/*
					Try to match a valid numeric value in one of the following formats:
					 - 123456
					 - 123.456
					 - 1_23.456e3
					 - 123.456E3
					 - 123.456e+3
					 - 123.456E+3
					 - 123.456e-3
					 - 123.456E-3
					 - 12345e5
				*/
				var (
					dotFound       = false
					eFound         = false
					postESignFound = false
					underscoreOk   = false
					number         string
				)

				for !t.isEnd() {
					if t.current() == '_' {
						if !underscoreOk || t.peek(-1) == '.' {
							t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalNumberFormat, t.line, t.col, ", unexpected or illegally placed underscore")
							goto start
						}
						t.advance()
					} else if t.current() == '.' {
						if !isNumber(t.peek(1)) {
							break
						}
						number += string(t.current())
						if dotFound {
							t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalNumberFormat, t.line, t.col, ", double decimal")
							goto start
						}
						dotFound = true
						t.advance()
						continue
					} else if isMatch('e', t.current()) {
						underscoreOk = false
						number += string(t.current())
						char := t.peek(1)

						if t.peekEnd(1) {
							t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalNumberFormat, t.line, t.col, ", missing exponent prefix")
							goto start
						} else if char != '+' && char != '-' && !isNumber(char) {
							t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalNumberFormat, t.line, t.col, ", expected '+', '-', or a digit")
							goto start
						}
						eFound = true
						t.advance()
						continue
					} else if eFound && isSign(t.current()) {
						number += string(t.current())
						if postESignFound {
							t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalNumberFormat, t.line, t.col, ", duplicate exponent sign postfix")
							goto start
						}
						postESignFound = true
						t.advance()
						continue
					} else if t.current() != '.' && !isNumber(t.current()) {
						break
					} else {
						if isNumber(t.current()) && !eFound {
							underscoreOk = true
						}
						number += string(t.current())
						t.advance()
					}
				}

				t.entities = append(t.entities, NewDefaultTokenEntity(t.line, t.col, IntegerLiteral, number))
				goto start
			}
		} else if t.current() == '"' { // start of a string literal
			value := ""
			if t.tokensLeft() < 2 {
				t.errManager.CreateNewErrorWithLineAndCol(ErrExpectedStrLiteralEOF, t.line, t.col, "")
				t.advance()
				goto start
			}
			t.advance()

			escaped := false
			escapeFound := false

			for !t.isEnd() {
				if t.current() == '\n' {
					t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalStrFormat, t.line, t.col, ", expected '\"' before end of line")
					t.newLine()
					goto start
				} else if !escaped && t.current() == '\\' {
					value += string(t.current())
					escaped = true
					escapeFound = true
					t.advance()
					continue
				} else if !escaped {
					if t.current() == '"' {
						break
					}
					value += string(t.current())
					escaped = false
				}

				t.advance()
			}

			if t.isEnd() {
				t.errManager.CreateNewErrorWithLineAndCol(ErrUnexpectedEOF, t.line, t.col, "")
				goto start
			}

			if !escapeFound {
				t.entities = append(t.entities, NewDefaultTokenEntity(t.line, t.col, StringLiteral, value))
			} else {
				t.entities = append(t.entities, NewDefaultTokenEntity(t.line, t.col, StringLiteral, getEscapedStr(value)))
			}

			t.advance()
			goto start
		} else if t.current() == '\'' { // start of a char
			var char string
			if t.tokensLeft() < 2 {
				t.errManager.CreateNewErrorWithLineAndCol(ErrExpectedCharLiteralEOF, t.line, t.col, "")
				t.advance()
				goto start
			}
			t.advance()

			escaped := false
			escapedFound := false
			hasChar := false

			for !t.isEnd() {
				if !escaped && t.current() == '\\' {
					if hasChar {
						t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalCharFormat, t.line, t.col, ", a char literal cannot contain more than one character")
						goto start
					}
					char += string(t.current())
					escaped = true
					escapedFound = true
					t.advance()
					continue
				} else if !escaped {
					if t.current() == '\\' {
						break
					}
					if hasChar {
						t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalCharFormat, t.line, t.col, ", a char literal cannot contain more than one character")
						goto start
					}

					hasChar = true
					char += string(t.current())
				} else if escaped {
					hasChar = true
					if !isLetter(byte(unicode.ToLower(rune(t.current())))) && t.current() == '\\' {
						t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalCharFormat, t.line, t.col, ", text preceding '\\' must be alphabetical")
						goto start
					}
					char += string(t.current())
					escaped = false
				}
				t.advance()
			}

			if t.isEnd() {
				t.errManager.CreateNewErrorWithLineAndCol(ErrUnexpectedEOF, t.line, t.col, "")
				goto start
			}

			if !escapedFound {
				if char == "" {
					t.errManager.CreateNewErrorWithLineAndCol(ErrIllegalCharFormat, t.line, t.col, ", char literal cannot be empty")
					goto start
				} else {
					t.entities = append(t.entities, NewDefaultTokenEntity(t.line, t.col, CharLiteral, char))
				}
			} else {
				t.entities = append(t.entities, NewDefaultTokenEntity(t.line, t.col, CharLiteral, getEscapedStr(char)))
			}

			t.advance()
			goto start
		} else {
			t.errManager.CreateNewErrorWithLineAndCol(ErrUnexpectedSymbol, t.line, t.col, ", '"+string(t.current())+"'")
			t.advance()
			goto start
		}
	}

end:
	t.entities = append(t.entities, t.EOFToken)
}

// parseLines parses lines.
func (t *Tokenizer) parseLines() {
	line := ""
	for i := 0; uint(i) < t.length; i++ {
		if t.tokens[i] == '\n' {
			t.lines = append(t.lines, line)
			line = ""
		} else {
			line += string(t.tokens[i])
		}
	}

	if line != "" {
		t.lines = append(t.lines, line)
	}
}

// Checks if the given character is a whitespace or not.
func isWhitespace(c byte) bool {
	return c == '\n' || c == ' ' || c == '\r' || c == '\t' || c == '\b' || c == '\f' || c == '\v'
}

// isEnd checks is the end is reached or not.
func (t *Tokenizer) isEnd() bool {
	return t.cursor >= t.length
}

// peekEnd checks if a jump forward for certain characters or not.
func (t *Tokenizer) peekEnd(f int) bool {
	return (int(t.cursor) + f) >= int(t.length)
}

// peek returns the character forward for certain characters.
func (t *Tokenizer) peek(f int) byte {
	if t.peekEnd(f) || (int(t.cursor)+f) < 0 {
		return t.tokens[t.length-1]
	} else {
		return t.tokens[int(t.cursor)+f]
	}
}

// current returns the current character.
func (t *Tokenizer) current() byte {
	if !t.isEnd() {
		return t.tokens[t.cursor]
	} else {
		return t.tokens[t.length-1]
	}
}

// isSymbol checks if the given character is a symbol or not.
func (t *Tokenizer) isSymbol(c byte) bool {
	return c == '+' || c == '-' ||
		c == '*' || c == '/' ||
		c == '^' || c == '<' ||
		c == '>' || c == '=' ||
		c == ',' || c == '(' ||
		c == ')' || c == '[' ||
		c == ']' || c == '{' ||
		c == '}' || c == '%' ||
		c == ':' || c == '?' ||
		c == '&' || c == '|' ||
		c == ';' || c == '!' ||
		c == '.' || c == '$' ||
		c == '#' || c == '@'
}

// isNumber checks if the given character is a digit or not.
func isNumber(c byte) bool {
	return unicode.IsDigit(rune(c))
}

// isLetter checks if the given character is a letter or not.
func isLetter(c byte) bool {
	return unicode.IsLetter(rune(c))
}

// isHexNum checks if the given character is a hex number or not.
func isHexNum(c byte) bool {
	return isNumber(c) || (c >= 65 && c <= 72) || (c >= 97 && c <= 104)
}

// advance advances the col and the cursor of the tokenizer.
func (t *Tokenizer) advance() {
	t.col++
	t.cursor++
}

// tokensLeft returns the number of tokens left.
func (t *Tokenizer) tokensLeft() uint {
	return t.length - t.cursor
}

// isSign checks if the given character is + or -.
func isSign(c byte) bool {
	return c == '+' || c == '-'
}

// newLine increments the line number and resets the column number.
func (t *Tokenizer) newLine() {
	t.col = 0
	t.line++
}

// commentStart checks if a comment is starting or not.
func commentStart(char0, char1 byte, mode *uint) bool {
	*mode = 0
	if char0 == '/' {
		if char1 == '/' {
			*mode = 1
		} else if char1 == '*' {
			*mode = 2
		}
	}
	return *mode != 0
}

// commentEnd checks if a comment is ending or not.
func commentEnd(char0, char1 byte, mode *uint) bool {
	return (*mode == 1 && char0 == '\n') || (*mode == 2 && char0 == '*' && char1 == '/')
}

// isMatch checks if two characters are the same or not regardless of capitalizing.
func isMatch(expected, char byte) bool {
	return unicode.ToLower(rune(expected)) == unicode.ToLower(rune(char))
}

// getEscapedStr formats the given string value to properly escaped string.
func getEscapedStr(unescaped string) string {
	value := ""
	for i := 0; i < len(unescaped); i++ {
		if value[i] == '\\' {
			switch unescaped[i+1] {
			case 'n':
				value += string('\n')
			case 't':
				value += string('\t')
			case 'b':
				value += string('\b')
			case 'v':
				value += string('\v')
			case 'r':
				value += string('\r')
			case 'f':
				value += string('\f')
			case 'a':
				value += string('\a')
			default:
				value += string(value[i+1])
			}
		}
	}
	return value
}

func (t *Tokenizer) GetData() string {
	return t.tokens
}

func (t *Tokenizer) GetErrors() *ErrorManager {
	return t.errManager
}
