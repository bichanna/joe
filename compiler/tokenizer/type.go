package tokenizer

const (
	Number      = iota // number
	Letter             // letter
	Underscore         // _
	LeftParen          // (
	LeftCurly          // {
	RightParen         // )
	RightCurly         // }
	Hash               // #
	Dot                // .
	Plus               // +
	Minus              // -
	Mult               // *
	Div                // /
	Mod                // %
	Colon              // :
	SemiColon          // ;
	DubQuote           // "
	SingQuote          // '
	Comma              // ,
	Newline            // \n
	LTE                // <=
	GTE                // >=
	EQEQ               // ==
	PlusEq             // +=
	MinusEq            // -=
	MultEq             // *=
	DivEq              // /=
	AndEq              // &=
	OrEq               // |=
	XOrEq              // ^=
	ModEq              // %=
	NotEq              // !=
	SHL                // <<
	SHR                // >>
	LessThan           // <
	GreaterThan        // >
	BitAnd             // &
	And                // &&
	BitOr              // |
	Or                 // ||
	XOr                // ^
	Not                // !
	Assign             // =
	LeftBrace          // [
	RightBrace         // ]
	Question           // ?
	PTR                // ->
	Inc                // ++
	Dec                // --
	Dollar             // $
	None               // not a type
	EOF         = 9999 // end of file
)
