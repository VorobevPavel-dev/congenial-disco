package table

import (
	"fmt"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

type Cell struct {
	Row    uint
	Column uint
	Value  string
}

func (c *Cell) AsString() string {
	return c.Value
}

func (c *Cell) AsInt() int64 {
	return 0
}

type LinearTable struct {
	Columns  []*parser.ColumnDefinition
	Elements *[][]Cell
	Name     *tokenizer.Token
}

func (lt LinearTable) IsInitialized() bool {
	return len(lt.Columns) > 0
}

// Create table will initialize columns inside fresh linear table
// It will take parser.ColumnDefinitions from request and append them to LinearTable.Columns slice
// Returns table name if all happened without errors
func (lt LinearTable) Create(req *parser.CreateTableQuery) (string, error) {
	lt.Columns = append(lt.Columns, req.Cols...)
	lt.Name = req.Name
	return lt.Name.Value, nil
}

func (lt LinearTable) Select(req *parser.SelectStatement) (*[][]Element, error) {
	return nil, nil
}

func (lt LinearTable) Insert(req *parser.InsertStatement) error {
	return nil
}

func (lt LinearTable) ShowCreate() string {
	initialRequest := fmt.Sprintf("CREATE TABLE %s (%s %s", lt.Name.Value, lt.Columns[0].Name.Value, strings.ToUpper(lt.Columns[0].Datatype.Value))
	for _, colDef := range lt.Columns[1:] {
		initialRequest += fmt.Sprintf(", %s %s", colDef.Name.Value, strings.ToUpper(colDef.Datatype.Value))
	}
	initialRequest += ");"
	return initialRequest
}

func (lt LinearTable) GetColumnsNames() []string {
	result := make([]string, len(lt.Columns))
	for index, colDef := range lt.Columns {
		result[index] = colDef.Name.Value
	}
	return result
}
