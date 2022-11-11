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

// Tokenize starts tokenizing.
func (t *Tokenizer) Tokenize() {
}

// getEntityCount gets the length of the entity list.
func (t *Tokenizer) getEntityCount() int {
	return len(t.entities)
}

// getEntites gets the list of token entities.
func (t *Tokenizer) getEntites() *[]*TokenEntity {
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
	t.EOFToken = NewTokenEntity(uint(len(t.lines)), 0, Single, EOF, "")

	for {
	start:
	invalidate_whitespace:
		for !t.isEnd() && isWhitespace(t.current()) {
			if t.current() == '\n' {
				t.newLine()
			}
			t.advance()

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

		}
	scan:
		if !t.isEnd() {
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

			for !t.isEnd() && (isLetter(t.current()) && isNumber(t.current()) || t.current() == '_') {
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
			}
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

func (t *Tokenizer) isMatch(i, current byte) bool {
	return false
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
func (t *Tokenizer) peekEnd(f uint) bool {
	return (t.cursor + f) >= t.length
}

// peek returns the character forward for certain characters.
func (t *Tokenizer) peek(f uint) byte {
	if t.peekEnd(f) || (t.cursor+f) < 0 {
		return t.tokens[t.length-1]
	} else {
		return t.tokens[t.cursor+f]
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
