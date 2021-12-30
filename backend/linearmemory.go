package backend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/VorobevPavel-dev/congenial-disco/ast"
	"strconv"
)

type MemoryCell []byte

func (mc MemoryCell) AsInt() int32 {
	var i int32
	err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i)
	if err != nil {
		panic(err)
	}
	return i
}

func (mc MemoryCell) AsText() string {
	return string(mc)
}

type table struct {
	columns     []string
	columnTypes []ColumnType
	rows        [][]MemoryCell
}

type MemoryBackend struct {
	tables map[string]*table
}

//func NewMemoryBackend() *MemoryBackend {
//	return &MemoryBackend{
//		tables: map[string]*table{},
//	}
//}

func (mb *MemoryBackend) tokenToCell(t *ast.Token) MemoryCell {
	if t.Kind == ast.NumericKind {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.Value)
		if err != nil {
			panic(err)
		}
		err = binary.Write(buf, binary.BigEndian, int32(i))
		if err != nil {
			panic(err)
		}
		return buf.Bytes()
	}
	if t.Kind == ast.StringKind {
		return MemoryCell(t.Value)
	}
	return nil
}

func (mb *MemoryBackend) CreateTable(crt *ast.CreateTableStatement) error {
	t := table{}
	mb.tables[crt.Name.Value] = &t
	columns := crt.Cols
	if columns == nil {
		return nil
	}
	for _, col := range *columns {
		t.columns = append(t.columns, col.Name.Value)
		var dt ColumnType
		switch col.Datatype.Value {
		case "int":
			dt = IntType
		case "text":
			dt = TextType
		default:
			return ErrInvalidDatatype
		}
		t.columnTypes = append(t.columnTypes, dt)
	}
	return nil
}

func (mb *MemoryBackend) Insert(inst *ast.InsertStatement) error {
	table, ok := mb.tables[inst.Table.Value]
	if !ok {
		return ErrTableDoesNotExist
	}
	insertValues := inst.Values
	if insertValues == nil {
		return nil
	}

	var row []MemoryCell
	if len(*insertValues) != len(table.columns) {
		return ErrMissingValues
	}

	for _, value := range *insertValues {
		if value.Kind != ast.LiteralKind {
			fmt.Println("Skipping non-literal")
			continue
		}
		row = append(row, mb.tokenToCell(value.Literal))
	}
	table.rows = append(table.rows, row)
	return nil
}

func (mb *MemoryBackend) Select(slct *ast.SelectStatement) (*Results, error) {
	table, ok := mb.tables[slct.From.Value]
	if !ok {
		return nil, ErrTableDoesNotExist
	}
	var results [][]Cell
	var columns []struct {
		Type ColumnType
		Name string
	}

	for i, row := range table.rows {
		var result []Cell
		isFirstRow := i == 0

		for _, exp := range slct.Item {
			if exp.Kind != ast.LiteralKind {
				fmt.Print("Skipping non-literal expression")
				continue
			}
			lit := exp.Literal
			if lit.Kind == ast.IdentifierKind {
				found := false
				for i, tableCol := range table.columns {
					if tableCol == lit.Value {
						if isFirstRow {
							columns = append(columns, struct {
								Type ColumnType
								Name string
							}{
								Type: table.columnTypes[i],
								Name: lit.Value,
							})
						}
						result = append(result, row[i])
						found = true
						break
					}
				}
				if !found {
					return nil, ErrColumnDoesNotExist
				}
				continue
			}
			return nil, ErrColumnDoesNotExist
		}
		results = append(results, result)
	}
	return &Results{
		Columns: columns,
		Rows:    results,
	}, nil
}
