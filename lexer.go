package main

import (
	"strings"
	"unicode"
)

type Lexer struct {
	source     string
	index      int
	filename   string
	line       int
	col        int
	tokens     []*Token
	errManager *ErrorManager
}

// newLexer creates a new lexer and returns it.
func newLexer(sourceStr string, filename string, aggressive bool) *Lexer {
	return &Lexer{
		source:     sourceStr,
		index:      0,
		filename:   filename,
		line:       1,
		col:        0,
		errManager: newErrorManager(filename, aggressive),
	}
}

// addToken appends the given token to the tokens list.
func (l *Lexer) addToken(token *Token) {
	l.tokens = append(l.tokens, token)
}

// currentPos returns a Position.
func (l *Lexer) currentPos() *Position {
	return &Position{
		file: l.filename,
		line: l.line,
		col:  l.col,
	}
}

// isEOF checks if it's at the end.
func (l *Lexer) isEOF() bool {
	return l.index == len(l.source)
}

// peek returns the current character.
func (l *Lexer) peek() rune {
	return rune(l.source[l.index])
}

// peekAhead returns the character at n characters ahead.
func (l *Lexer) peekAhead(n int) rune {
	if l.index+n >= len(l.source) {
		return ' '
	}
	return rune(l.source[l.index+n])
}

// next returns the current character and moves the index one ahead.
func (l *Lexer) next() rune {
	char := l.source[l.index]

	if l.index < len(l.source) {
		l.index++
	}

	if char == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}

	return rune(char)
}

// back moves the index back one.
func (l *Lexer) back() {
	if l.index > 0 {
		l.index--
		if l.source[l.index] == '\n' {
			l.line--
			count := -1
			i := l.index - 1
			for l.source[i] != '\n' {
				count++
				i++
			}
		} else {
			l.col--
		}
	}
}

// readUntil reads characters until encountering char.
func (l *Lexer) readUntil(char rune) string {
	value := strings.Builder{}
	for !l.isEOF() && l.peek() != char {
		value.WriteRune(l.next())
	}
	return value.String()
}

// readIdentifier reads an identifier and returns it.
func (l *Lexer) readIdentifier() string {
	identifier := strings.Builder{}
	for {
		if l.isEOF() {
			break
		}

		c := l.next()
		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '?' || c == '!' || c == '-' {
			identifier.WriteRune(c)
		} else {
			l.back()
			break
		}
	}

	return identifier.String()
}

// readNumeral reads a numeral and returns the number and whether it is a floating point number or not.
func (l *Lexer) readNumeral() (string, bool) {
	sawDot := false
	number := strings.Builder{}
	for {
		if l.isEOF() {
			break
		}

		c := l.next()
		if unicode.IsDigit(c) {
			number.WriteRune(c)
		} else if c == '.' && !sawDot {
			sawDot = true
			number.WriteRune(c)
		} else if c == '.' && sawDot {
			l.errManager.newErrorWithPosition(l.currentPos(), "invalid floating point number format")
			break
		} else {
			l.back()
			break
		}
	}
	return number.String(), sawDot
}

// readString reads a string and returns the value.
func (l *Lexer) readString() string {
	l.next()
	value := strings.Builder{}
	for !l.isEOF() && l.peek() != '\'' && l.peek() != '"' {
		char := l.next()
		if char == '\\' {
			if l.isEOF() {
				l.errManager.newErrorWithPosition(l.currentPos(), "string literal not enclosed with '\"'")
				break
			} else {
				char = l.next()
				switch char {
				case 'a':
					value.WriteRune('\a')
				case 'b':
					value.WriteRune('\b')
				case 'r':
					value.WriteRune('\r')
				case '\\':
					value.WriteRune('\\')
				case 't':
					value.WriteRune('\t')
				case 'n':
					value.WriteRune('\n')
				default:
					value.WriteRune(char)
				}
				continue
			}
		}
		value.WriteRune(char)
	}

	// read the ending quote
	l.next()

	return value.String()
}

// tokenizeNextToken tokenize a token.
func (l *Lexer) tokenizeNextToken() {
	char := l.next()

	switch char {
	case ' ', '\t', '\n':
		// do nothing
		break
	case '#':
		l.readUntil('\n')
		// do nothing
		break
	case '(':
		l.addToken(newDefaultToken(leftParen, *l.currentPos()))
	case ')':
		l.addToken(newDefaultToken(rightParen, *l.currentPos()))
	case '[':
		l.addToken(newDefaultToken(leftBracket, *l.currentPos()))
	case ']':
		l.addToken(newDefaultToken(rightBracket, *l.currentPos()))
	case '{':
		l.addToken(newDefaultToken(leftBrace, *l.currentPos()))
	case '}':
		l.addToken(newDefaultToken(rightBrace, *l.currentPos()))
	case '.':
		l.addToken(newDefaultToken(dot, *l.currentPos()))
	case '+':
		l.addToken(newDefaultToken(plus, *l.currentPos()))
	case '-':
		l.addToken(newDefaultToken(minus, *l.currentPos()))
	case '*':
		l.addToken(newDefaultToken(times, *l.currentPos()))
	case '/':
		l.addToken(newDefaultToken(div, *l.currentPos()))
	case '%':
		l.addToken(newDefaultToken(mod, *l.currentPos()))
	case '>':
		if l.peekAhead(1) != '=' {
			l.addToken(newDefaultToken(greater, *l.currentPos()))
		} else {
			pos := *l.currentPos()
			l.next()
			l.addToken(newDefaultToken(geq, pos))
		}
	case '<':
		if l.peekAhead(1) != '=' {
			l.addToken(newDefaultToken(less, *l.currentPos()))
		} else {
			pos := *l.currentPos()
			l.next()
			l.addToken(newDefaultToken(leq, pos))
		}
	case '!':
		if l.peekAhead(1) != '=' {
			l.errManager.newErrorWithPosition(l.currentPos(), "expected '='")
			l.addToken(newDefaultToken(unknown, *l.currentPos()))
		} else {
			pos := *l.currentPos()
			l.next()
			l.addToken(newDefaultToken(neq, pos))
		}
	case '=':
		l.addToken(newDefaultToken(eq, *l.currentPos()))
	case '\'', '"':
		l.addToken(newToken(strLiteral, *l.currentPos(), l.readString()))
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		pos := l.currentPos()
		number, isDouble := l.readNumeral()
		number = string(char) + number
		if !isDouble {
			l.addToken(newToken(intLiteral, *pos, number))
		} else {
			l.addToken(newToken(doubleLiteral, *pos, number))
		}
	default:
		pos := l.currentPos()
		payload := string(char) + l.readIdentifier()
		switch payload {
		case "and":
			l.addToken(newDefaultToken(and, *pos))
		case "or":
			l.addToken(newDefaultToken(or, *pos))
		case "not":
			l.addToken(newDefaultToken(not, *pos))
		case "if":
			l.addToken(newDefaultToken(ifKW, *pos))
		case "while":
			l.addToken(newDefaultToken(whileKW, *pos))
		case "for":
			l.addToken(newDefaultToken(forKW, *pos))
		case "mut":
			l.addToken(newDefaultToken(mutKW, *pos))
		case "let":
			l.addToken(newDefaultToken(letKW, *pos))
		case "set":
			l.addToken(newDefaultToken(setKW, *pos))
		case "map":
			l.addToken(newDefaultToken(mapKW, *pos))
		case "def":
			l.addToken(newDefaultToken(defineKW, *pos))
		case "lambda":
			l.addToken(newDefaultToken(lambdaKW, *pos))
		case "block":
			l.addToken(newDefaultToken(blockKW, *pos))
		case "async":
			l.addToken(newDefaultToken(asyncKW, *pos))
		case "await":
			l.addToken(newDefaultToken(awaitKW, *pos))
		case "struct":
			l.addToken(newDefaultToken(structKW, *pos))
		case "true":
			l.addToken(newDefaultToken(trueLiteral, *pos))
		case "false":
			l.addToken(newDefaultToken(falseLiteral, *pos))
		case "nil":
			l.addToken(newDefaultToken(nilLiteral, *pos))
		default:
			l.addToken(newToken(identifier, *pos, payload))
		}
	}
}

// tokenize tokenizes a source file and returns the tokens.
func (l *Lexer) tokenize() []*Token {
	for !l.isEOF() {
		l.tokenizeNextToken()
	}

	l.addToken(newDefaultToken(EOF, *l.currentPos()))

	return l.tokens
}
