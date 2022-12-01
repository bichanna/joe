package main

import "fmt"

type Position struct {
	file string
	line int
	col  int
}

// String returns stringified Position.
func (p *Position) String() string {
	return fmt.Sprintf("[%d:%d]", p.line, p.col)
}

type Token struct {
	kind    tokenType
	pos     Position
	payload string
}

// newToken creates a new token.
func newToken(kind tokenType, pos Position, payload string) *Token {
	return &Token{
		kind:    kind,
		pos:     pos,
		payload: payload,
	}
}

func newDefaultToken(kind tokenType, pos Position) *Token {
	return &Token{
		kind: kind,
		pos:  pos,
	}
}

// String stringifies the token.
func (t *Token) String() string {
	if t.kind == identifier {
		return t.payload
	} else {
		switch t.kind {
		case plus:
			return "+"
		case minus:
			return "-"
		case times:
			return "*"
		case div:
			return "/"
		case mod:
			return "%"
		case greater:
			return ">"
		case less:
			return "<"
		case eq:
			return "="
		case geq:
			return ">="
		case leq:
			return "<="
		case neq:
			return "!="
		case and:
			return "&&"
		case or:
			return "||"
		case not:
			return "!"
		case ifKW:
			return "if"
		case whileKW:
			return "while"
		case forKW:
			return "for"
		case mutKW:
			return "mut"
		case letKW:
			return "let"
		case setKW:
			return "set"
		case mapKW:
			return "map"
		case listKW:
			return "list"
		case defineKW:
			return "def"
		case lambdaKW:
			return "lambda"
		case blockKW:
			return "block"
		case structKW:
			return "struct"
		case asyncKW:
			return "async"
		case awaitKW:
			return "await"
		case trueLiteral:
			return "true"
		case falseLiteral:
			return "false"
		case nilLiteral:
			return "nil"
		default:
			return "UNKNOWN"
		}
	}
}

type tokenType uint8

const (
	unknown tokenType = iota
	leftParen
	rightParen
	leftBracket
	rightBracket
	leftBrace
	rightBrace
	dot
	identifier

	// binary operators
	plus
	minus
	times
	div
	mod
	greater
	less
	eq
	geq
	leq
	neq
	and
	or
	not

	// keywords
	ifKW
	whileKW
	forKW
	mutKW
	letKW
	setKW
	mapKW
	listKW
	defineKW
	lambdaKW
	blockKW
	structKW
	asyncKW
	awaitKW

	// literals
	trueLiteral
	falseLiteral
	nilLiteral
	strLiteral
	intLiteral
	doubleLiteral

	EOF
)
