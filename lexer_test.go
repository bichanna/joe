package main

import (
	"testing"
)

func TestNonAggressiveLexer(t *testing.T) {
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
