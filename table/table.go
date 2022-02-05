package table

import (
	"github.com/VorobevPavel-dev/congenial-disco/parser"
)

type Element interface {
	AsString() string
	AsInt() int64
}

type Table interface {
	Create(req *parser.CreateTableStatement) error
	Select(req *parser.SelectStatement) (*[][]Element, error)
	Insert(req *parser.InsertStatement) error
}
