package table

import (
	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

type Table interface {
	IsInitialized() bool
	Create(req *parser.CreateTableQuery) (Table, string, error)
	Select(req *parser.SelectQuery) ([][]*tokenizer.Token, error)
	Insert(req *parser.InsertIntoQuery) (Table, error)
	Engine() string
	ShowCreate() string
	GetColumns() string
	Count() int
}
