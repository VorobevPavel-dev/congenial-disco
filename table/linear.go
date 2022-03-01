package table

import (
	"fmt"
	"strings"

	p "github.com/VorobevPavel-dev/congenial-disco/parser"
	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

//LinearTable is a struct allowing users to store data
// in PostgreSQL-like format - by rows. Elements inside will be stored in a
// two-dimensional slice of parser.element structures
type LinearTable struct {
	//columns will contain column names with indexes from CREATE TABLE request
	columns     map[int]string
	elements    [][]element
	name        string
	columnTypes []t.TokenKind
}

func (lt LinearTable) IsInitialized() bool {
	// len of empty map is 0
	return len(lt.columns) > 0
}

// Getters

func (lt LinearTable) Engine() string {
	return "linear"
}

func (lt LinearTable) Name() string {
	return lt.name
}

func (lt LinearTable) ShowCreate() string {
	return ""
}

// Handlers

func (lt LinearTable) Create(req *p.CreateTableQuery) (Table, string, error) {
	lt.columns = make(map[int]string)
	lt.columnTypes = make([]t.TokenKind, len(*req.Cols))
	for i, definition := range *req.Cols {
		lt.columns[i] = definition.Name.Value
		lt.columnTypes[i] = typeCorreation[definition.Datatype.Value]
	}
	lt.name = req.Name.Value
	return Table(lt), lt.name, nil
}

func (lt LinearTable) Insert(req *p.InsertIntoQuery) (Table, error) {
	var valuesToInsert [][]element
	for _, set := range req.Values {
		// Validate set has right lenght
		if len(set) != len(lt.columns) {
			return nil, fmt.Errorf(ErrIncorrectSetLenghtTemplate, len(lt.columns), len(set))
		}
		// Consert set to slice of elements
		tempElementSet := make([]element, len(set))
		for i, token := range set {
			tempElementSet[i] = tokenToElement(token)
		}
		// Assert elements has correct types
		for elementIndex, element := range tempElementSet {
			if element.kind != lt.columnTypes[elementIndex] {
				return nil, fmt.Errorf(ErrIncorrectValueType, element.value, element.kind)
			}
		}
		valuesToInsert = append(valuesToInsert, tempElementSet)
	}
	// Insert sets to table
	lt.elements = append(lt.elements, valuesToInsert...)
	return Table(lt), nil
}

func (lt LinearTable) Select(req *p.SelectQuery) ([][]string, error) {
	// colIndex will allow to get column index number by its name
	var colIndex = func(t string) int {
		for i, n := range lt.columns {
			if strings.Compare(n, t) == 0 {
				return i
			}
		}
		return -1
	}
	// Get column indexes to select from table
	columnIndexesToSelect := []int{}
	for _, columnToSelect := range req.Columns {
		if index := colIndex(columnToSelect.Value); index != -1 {
			columnIndexesToSelect = append(columnIndexesToSelect, index)
		} else {
			return [][]string{}, fmt.Errorf(ErrIncorrectSelectColumn, columnToSelect.Value, lt.name)
		}
	}
	// Get conditions for every column in table by its index
	conditionCorrelation := make(map[int][]*p.Condition)
	for _, condition := range req.Conditions {
		if index := colIndex(condition.Column.Value); index == -1 {
			return [][]string{}, fmt.Errorf(ErrIncorrectConditionColumn, condition.Column.Value, lt.name)
		} else {
			conditionCorrelation[index] = append(conditionCorrelation[index], condition)
		}
	}
	// Construct header for output
	header := make([]string, len(columnIndexesToSelect))
	for i, toSelect := range columnIndexesToSelect {
		header[i] = lt.columns[toSelect]
	}
	result := [][]string{header}

	// Go around the table
	for _, currentRow := range lt.elements {
		var selectRow bool = true
		for columnIndex, element := range currentRow {
			for _, condition := range conditionCorrelation[columnIndex] {
				passes, err := condition.EvaluateWithValues(&t.Token{Value: element.value, Kind: element.kind})
				if err != nil {
					return nil, fmt.Errorf(ErrConditionRuntime, err)
				}
				if !passes {
					selectRow = false
				}
			}
		}

		// Get necessary values if row passes all checks
		extractedValues := []string{}
		if selectRow {
			for _, columnIndex := range columnIndexesToSelect {
				extractedValues = append(extractedValues, currentRow[columnIndex].value)
			}
			result = append(result, extractedValues)
		}
	}
	return result, nil
}
