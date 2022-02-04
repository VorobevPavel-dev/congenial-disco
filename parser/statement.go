package parser

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
}

func Parse(request string) *Statement {
	return nil
}
