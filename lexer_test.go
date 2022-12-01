package main

import (
	"testing"
)

func Test1(t *testing.T) {
	source := `(println "Hello World")
	# Hello World
	(def is-even? [num] (= 0 (% num 2)))`

	lexer := newLexer(source, "[input]", false)
	tokens := lexer.tokenize()

	if len(tokens) != 21 {
		// It's 21 because EOF token is also included.
		t.Error("Invalid, got ", len(tokens))
	}
}

func Test2(t *testing.T) {
	source := `"Hello \n World"`

	lexer := newLexer(source, "[input]", false)
	tokens := lexer.tokenize()

	if len(tokens) != 2 {
		t.Error("Invalid, got", len(tokens))
	}

	if tokens[0].payload == "Hello \n World" {
		t.Error("Invalid got \"", tokens[0].payload, "\"")
	}
}
