package parser

import "github.com/VorobevPavel-dev/congenial-disco/tokenizer"

type ColumnDefinition struct {
	Name     tokenizer.Token
	Datatype tokenizer.Token
}

type CreateTableStatement struct {
	Name tokenizer.Token
	Cols *[]*ColumnDefinition
}
