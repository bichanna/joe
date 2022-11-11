package compiler

type TokenEntityID uint

const (
	Identifier TokenEntityID = iota
	NativeType
	TypeIdentifier
	IntegerLiteral
	HexLiteral
	ModuleName
	StringLiteral
	CharLiteral
	Literal
	Value
	AccessType
	SingleLineComment
	ReturnStatement

	Single
	NoEntity
)

type TokenEntity struct {
	Line  uint
	Col   uint
	id    TokenEntityID
	ttype TokenType
	value string
}

// NewTokenEntity creates a new token entity.
func NewTokenEntity(line, col uint, id TokenEntityID, tt TokenType, value string) *TokenEntity {
	return &TokenEntity{
		Line:  line,
		Col:   col,
		id:    id,
		ttype: tt,
		value: value,
	}
}

// NewDefaultTokenEntity creates a new default (None) token entity.
func NewDefaultTokenEntity(line, col uint, id TokenEntityID, value string) *TokenEntity {
	return NewTokenEntity(line, col, id, None, value)
}

// GetID returns the TokenEntityID of the token entity.
func (te *TokenEntity) GetID() TokenEntityID {
	return te.id
}

// SetID sets a new TokenEntityID to the token entity.
func (te *TokenEntity) SetID(id TokenEntityID) {
	te.id = id
}

// GetValue returns the value of the token entity.
func (te *TokenEntity) GetValue() string {
	return te.value
}

// IsSingle checks if the token type is Single or not.
func (te *TokenEntity) IsSingle() bool {
	return te.id == Single
}

// GetTokenType returns the TokenType of the token entity.
func (te *TokenEntity) GetTokenType() TokenType {
	return te.ttype
}
