package ast

type Kind uint
type ExpressionKind uint

//const (
//	SelectKind AstKind = iota
//	CreateTableKind
//	InsertKind
//)

const (
	LiteralKind ExpressionKind = iota
)

type Expression struct {
	Literal *Token
	Kind    ExpressionKind
}

// Definition of structures for Insert statements

type InsertStatement struct {
	Table  Token
	Values *[]*Expression
}

// Definition of structures for Create statements

type ColumnDefinition struct {
	Name     Token
	Datatype Token
}

type CreateTableStatement struct {
	Name Token
	Cols *[]*ColumnDefinition
}

// Definition of structures for Select statements

type SelectStatement struct {
	Item []*Expression
	From Token
}
