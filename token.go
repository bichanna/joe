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
