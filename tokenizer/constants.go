package tokenizer

type TokenKind uint

const (
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
	// EngineKind will determine type of table
	EngineKind

	// charset is needed for generating random tokens
	charset string = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Constants is a constructor returning all reserver words and
// symbols that can be parsed into tokens.
func Constants() map[TokenKind]map[string]*Token {
	return map[TokenKind]map[string]*Token{
		KeywordKind: {
			"select":   {Value: "select", Kind: KeywordKind},
			"from":     {Value: "from", Kind: KeywordKind},
			"as":       {Value: "as", Kind: KeywordKind},
			"table":    {Value: "table", Kind: KeywordKind},
			"create":   {Value: "create", Kind: KeywordKind},
			"insert":   {Value: "insert", Kind: KeywordKind},
			"into":     {Value: "into", Kind: KeywordKind},
			"values":   {Value: "values", Kind: KeywordKind},
			"show":     {Value: "show", Kind: KeywordKind},
			"with":     {Value: "with", Kind: KeywordKind},
			"engine":   {Value: "engine", Kind: KeywordKind},
			"settings": {Value: "settings", Kind: KeywordKind},
			"where":    {Value: "where", Kind: KeywordKind},
			"order":    {Value: "order", Kind: KeywordKind},
			"by":       {Value: "by", Kind: KeywordKind},
		},
		SymbolKind: {
			";":  {Value: ";", Kind: SymbolKind},
			"*":  {Value: "*", Kind: SymbolKind},
			",":  {Value: ",", Kind: SymbolKind},
			"(":  {Value: "(", Kind: SymbolKind},
			")":  {Value: ")", Kind: SymbolKind},
			" ":  {Value: " ", Kind: SymbolKind},
			"==": {Value: "==", Kind: SymbolKind},
			">":  {Value: ">", Kind: SymbolKind},
			"<":  {Value: "<", Kind: SymbolKind},
		},
		TypeKind: {
			"int":  {Value: "int", Kind: TypeKind},
			"text": {Value: "text", Kind: TypeKind},
		},
		EngineKind: {
			"linear": {Value: "linear", Kind: EngineKind},
			"column": {Value: "column", Kind: EngineKind},
			// "blocked": {Value: "blocked", Kind: EngineKind},
		},
	}
}

// KindMap will return map where all kinds in keys has string representation in values
func KindMap() *map[TokenKind]string {
	return &map[TokenKind]string{
		NumericKind:    "number",
		KeywordKind:    "keyword",
		SymbolKind:     "symbol",
		IdentifierKind: "identifier",
		TypeKind:       "type",
		EngineKind:     "engine",
	}
}

// Keys will return only keys of first-level from given map
func Keys(input map[string]*Token) *[]string {
	result := make([]string, len(input))
	index := 0
	for key := range input {
		result[index] = key
		index++
	}
	return &result
}
