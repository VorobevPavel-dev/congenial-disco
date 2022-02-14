package tokenizer

type TokenKind uint

func keywordsMap() *map[string]Token {
	return &map[string]Token{
		"select": {Value: "select", Kind: KeywordKind},
		"from":   {Value: "from", Kind: KeywordKind},
		"as":     {Value: "as", Kind: KeywordKind},
		"table":  {Value: "table", Kind: KeywordKind},
		"create": {Value: "create", Kind: KeywordKind},
		"insert": {Value: "insert", Kind: KeywordKind},
		"into":   {Value: "into", Kind: KeywordKind},
		"values": {Value: "values", Kind: KeywordKind},
		"show":   {Value: "show", Kind: KeywordKind},
		"tables": {Value: "tables", Kind: KeywordKind},
	}
}

func symbolsMap() *map[string]Token {
	return &map[string]Token{
		";": {Value: ";", Kind: SymbolKind},
		"*": {Value: "*", Kind: SymbolKind},
		",": {Value: ",", Kind: SymbolKind},
		"(": {Value: "(", Kind: SymbolKind},
		")": {Value: ")", Kind: SymbolKind},
		" ": {Value: " ", Kind: SymbolKind},
	}
}

func typesMap() *map[string]Token {
	return &map[string]Token{
		"int":  {Value: "int", Kind: TypeKind},
		"text": {Value: "text", Kind: TypeKind},
	}
}

// constants builds a map where all tokens sorted by its Kind
func constants() *map[TokenKind](*map[string]Token) {
	return &map[TokenKind]*map[string]Token{
		KeywordKind: keywordsMap(),
		SymbolKind:  symbolsMap(),
		TypeKind:    typesMap(),
	}
}

const (
	// Charset for generating different tokens. Needed for tests
	charset string = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// NumericKind will correspond to all numeric values
	NumericKind TokenKind = iota
	// KeywordKind will correspond to all string equal to one of keywords
	KeywordKind
	// SymbolKind will correspond to every specified utility symbol
	SymbolKind
	// IdentifierKind will correspond to every custom value (table name, column name, values etc...)
	IdentifierKind
	// TypeKind will correspond to every type in request
	TypeKind
)

// KindToString returns string representation of token kind
func KindToString(kind int) string {
	kindCorrelation := map[TokenKind]string{
		TypeKind:       "type",
		NumericKind:    "number",
		KeywordKind:    "keyword",
		SymbolKind:     "symbol",
		IdentifierKind: "identifier",
	}
	if value, ok := kindCorrelation[TokenKind(kind)]; ok {
		return value
	}
	return "unknown"
}
