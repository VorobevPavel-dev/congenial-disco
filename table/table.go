package table

import (
	"github.com/VorobevPavel-dev/congenial-disco/parser"
)

type Element interface {
	AsString() string
	AsInt() (int64, error)
}

type Table interface {
	IsInitialized() bool
	Create(req *parser.CreateTableQuery) (string, error)
	Select(req *parser.SelectStatement) (*[][]Element, error)
	Insert(req *parser.InsertStatement) error
	ShowCreate() string
}
