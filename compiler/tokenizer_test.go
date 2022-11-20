package compiler

import "testing"

func TestTokenize(t *testing.T) {
	source := `print("Hello World\n");
	for i := 1; i < 100; i++ {
		print(array[i] + '\n');
	}`
	InitializeErrors()
	tokenizer := NewTokenizer(source, "test.joe")
	
	if tokenizer.GetEntityCount() != 25 {
		t.Error("the entity count is not 25!")
	}
}
