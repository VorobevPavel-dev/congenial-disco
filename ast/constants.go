package ast

//type keyword string
//type symbol string

type TokenKind uint

// SQL-reserved words
//const (
//	SelectKeyword keyword = "select"
//	FromKeyword   keyword = "from"
//	AsKeyword     keyword = "as"
//	TableKeyword  keyword = "table"
//	CreateKeyword keyword = "create"
//	InsertKeyword keyword = "insert"
//	IntoKeyword   keyword = "into"
//	ValuesKeyword keyword = "values"
//)

// Types keywords
//const (
//	IntKeyword  keyword = "int"
//	TextKeyword keyword = "text"
//)

// Symbol constants
//const (
//	SemicolonSymbol  symbol = ";"
//	AsteriskSymbol   symbol = "*"
//	CommaSymbol      symbol = ","
//	LeftParenSymbol  symbol = "("
//	RightParenSymbol symbol = ")"
//)

//KeywordKind TokenKind = iota
//SymbolKind

const (
	IdentifierKind TokenKind = iota
	StringKind
	NumericKind
)

type Location struct {
	Line uint
	Col  uint
}

type Token struct {
	Value string
	Kind  TokenKind
	Loc   Location
}

func (t *Token) equals(other *Token) bool {
	return t.Value == other.Value && t.Kind == other.Kind
}
