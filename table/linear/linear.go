package linear

import (
	"fmt"
	"sort"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table"
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
	"github.com/VorobevPavel-dev/congenial-disco/utility"
)

type LinearTable struct {
	Columns  []parser.ColumnDefinition
	Elements [][]*tokenizer.Token
	Name     *tokenizer.Token
}

func (lt LinearTable) IsInitialized() bool {
	return len(lt.Columns) > 0
}

// Create table will initialize columns inside fresh linear table
// It will take parser.ColumnDefinitions from request and append them to LinearTable.Columns slice
// Returns table name if all happened without errors
func (lt LinearTable) Create(req *parser.CreateTableQuery) (table.Table, string, error) {
	lt.Columns = make([]parser.ColumnDefinition, len(*req.Cols))
	for i := range *req.Cols {
		lt.Columns[i] = (*req.Cols)[i]
	}
	lt.Name = req.Name
	return table.Table(lt), lt.Name.Value, nil
}

func (lt LinearTable) Select(req *parser.SelectQuery) ([][]*tokenizer.Token, error) {
	tableColumns := lt.GetColumnsNames()
	// Build slice of indexes of columns to select
	indexes := []int{}
	for _, columnToSelect := range req.Columns {
		// Assert that all columns in request are actually inside table
		if !utility.StringIsIn(columnToSelect.Value, tableColumns) {
			return nil, fmt.Errorf("column with name %s is not represented inside table (actual columns: [%s])",
				columnToSelect,
				lt.GetColumns(),
			)
		}
		indexes = append(indexes, utility.FindStringInSlice(tableColumns, columnToSelect.Value))
	}

	// Sort indexes for pretty-print
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})

	header := make([]*tokenizer.Token, len(indexes))
	for i, extractedHeaderPosition := range indexes {
		header[i] = lt.Columns[extractedHeaderPosition].Name
	}

	// Go around table
	result := [][]*tokenizer.Token{}
	result = append(result, header)
	for _, currentRow := range lt.Elements {
		nullCounter := 0
		extractedValues := make([]*tokenizer.Token, len(lt.Columns))
		for i, columnIndex := range indexes {
			value := currentRow[columnIndex]
			if value.Value == "null" {
				nullCounter++
			}
			extractedValues[i] = currentRow[columnIndex]
			if nullCounter == len(indexes) {
				continue
			}
		}
		result = append(result, extractedValues)
	}
	return result, nil
}

func (lt LinearTable) Insert(req *parser.InsertIntoQuery) (table.Table, error) {
	// typeCorreation[column.Datatype.Value] will give Kind which is supported on current column
	typeCorreation := map[string]tokenizer.TokenKind{
		"int":  tokenizer.NumericKind,
		"text": tokenizer.IdentifierKind,
	}
	// Get kind of type for each column
	columnCorrelation := make(map[string]tokenizer.TokenKind)
	for _, column := range lt.Columns {
		columnCorrelation[column.Name.Value] = typeCorreation[column.Datatype.Value]
	}

	// Assert that all value sets has same size as len(lt.Columns)
	for _, valueSet := range req.Values {
		if len(valueSet) != len(lt.Columns) {
			return nil, fmt.Errorf(
				"value set ahs incorrect number of values: %s",
				tokenizer.Bracketize(valueSet),
			)
		}
	}

	// Construct map for insertion.
	// For example:
	//	- table has columns id, num, name
	//	- values are (1,2,test), (3,4, test2)
	// So mapToInsert will be
	// id: [1, 3]
	// num: [2, 4]
	// name: [test, test2]
	mapToInsert := make(map[string][]*tokenizer.Token)
	for columnIndex, column := range lt.Columns {
		intendedValues := []*tokenizer.Token{}
		for valueSetIndex := range req.Values {
			intendedValues = append(intendedValues, req.Values[valueSetIndex][columnIndex])
		}
		mapToInsert[column.Name.Value] = intendedValues
	}

	// Assert that all values for corresponding column have specified in table kind.
	// For example if table has columns id int, num int, name test
	// then all values from mapToInsert["id"] must be NumericType
	// all values from mapToInsert["num"] must be NumericType
	// all values from mapToInsert["name"] must be IdentifierType
	for columnName, insertingValues := range mapToInsert {
		expectedType := columnCorrelation[columnName]
		for i, value := range insertingValues {
			// get original set for error message
			failedIndex := []*tokenizer.Token{}
			for j := range req.Values {
				failedIndex = append(failedIndex, req.Values[j][i])
			}
			if value.Kind != expectedType {
				return nil, fmt.Errorf(
					"for valueset %s value %s has incorrect type (expected %s, got: %s",
					tokenizer.Bracketize(failedIndex),
					value.Value,
					tokenizer.KindToString(value.Kind),
					tokenizer.KindToString(expectedType),
				)
			}
		}
	}
	// If everything ok - get original values and insert them inside
	lt.Elements = append(lt.Elements, req.Values...)
	return lt, nil
}

func (lt LinearTable) ShowCreate() string {
	initialRequest := fmt.Sprintf("CREATE TABLE %s (%s %s", lt.Name.Value, lt.Columns[0].Name.Value, strings.ToUpper(lt.Columns[0].Datatype.Value))
	for _, colDef := range lt.Columns[1:] {
		initialRequest += fmt.Sprintf(", %s %s", colDef.Name.Value, strings.ToUpper(colDef.Datatype.Value))
	}
	initialRequest += ");"
	return initialRequest
}

func (lt LinearTable) Engine() string {
	return "Linear"
}

func (lt LinearTable) GetColumnsNames() []string {
	result := make([]string, len(lt.Columns))
	for index, colDef := range lt.Columns {
		result[index] = colDef.Name.Value
	}
	return result
}

func (lt LinearTable) GetColumns() string {
	var result string
	if !lt.IsInitialized() {
		return ""
	}
	for _, column := range lt.Columns {
		result += fmt.Sprintf("(%s, %s), ", column.Name.Value, column.Datatype.Value)
	}
	return result[:len(result)-2]
}

func (lt LinearTable) Count() int {
	return len(lt.Elements)
}
