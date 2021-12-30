package backend

import (
	"errors"
	"github.com/VorobevPavel-dev/congenial-disco/ast"
)

type ColumnType uint

const (
	TextType ColumnType = iota
	IntType
)

var (
	ErrInvalidDatatype    = errors.New("invalid datatype")
	ErrTableDoesNotExist  = errors.New("specified table does not exist")
	ErrMissingValues      = errors.New("not enough values provided for request")
	ErrColumnDoesNotExist = errors.New("provided column does not exist")
)

type Cell interface {
	AsText() string
	AsInt() int32
}

type Results struct {
	Columns []struct {
		Type ColumnType
		Name string
	}
	Rows [][]Cell
}

type Backend interface {
	CreateTable(statement *ast.CreateTableStatement) error
	Insert(statement *ast.InsertStatement) error
	Select(statement *ast.SelectStatement) error
}
