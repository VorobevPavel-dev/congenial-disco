package linear

import (
	"fmt"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table"
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
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

	// Map with correlation between column index and column name
	// For example if there is a table with columns id int, name text
	// indexCorrelation will be
	// id: 0
	// name: 1
	indexCorrelation := make(map[string]int)
	for index, column := range tableColumns {
		indexCorrelation[column] = index
	}

	indexesToSelect := []int{}

	// Get indexes of columns to select
	for _, columnToSelect := range req.Columns {
		if index, ok := indexCorrelation[columnToSelect.Value]; !ok {
			return nil, fmt.Errorf("column with name %s is not represented inside table (actual columns: [%s])",
				columnToSelect,
				lt.GetColumns(),
			)
		} else {
			indexesToSelect = append(indexesToSelect, index)
		}
	}

	// Map with correlation between column index and array of conditions on it
	// For table with columns id int, name text
	// and conditions id > 5 and id <10 and name == test
	// conditionCorrelation will be
	// 0: [{id > 5}, {id < 10}]
	// 1: [{name == test}]
	var (
		index int
		ok    bool
	)
	conditionCorrelation := make(map[int][]*parser.Condition)
	for _, condition := range req.Conditions {
		// Get index of column under condition
		if index, ok = indexCorrelation[condition.Column.Value]; !ok {
			return nil, fmt.Errorf("column %s not in table", condition.Column.Value)
		}
		conditionCorrelation[index] = append(conditionCorrelation[index], condition)
	}

	// Construct header for output. Contains column names
	header := make([]*tokenizer.Token, len(indexesToSelect))
	for i, extractedHeaderPosition := range indexesToSelect {
		header[i] = lt.Columns[extractedHeaderPosition].Name
	}

	// Go around table
	result := [][]*tokenizer.Token{header}

	var selectRow bool = true
	for _, currentRow := range lt.Elements {
		// Check if row passes all conditions
		for columnIndex, element := range currentRow {
			for _, condition := range conditionCorrelation[columnIndex] {
				passes, err := condition.EvaluateWithValues(element)
				if err != nil {
					return nil, fmt.Errorf("error while checking conditions: %v", err)
				}
				if !passes {
					selectRow = false
					break
				}
			}
		}
		// Get necessary columns
		extractedValues := []*tokenizer.Token{}
		if selectRow {
			for _, columnIndex := range indexesToSelect {
				extractedValues = append(extractedValues, currentRow[columnIndex])
			}
			result = append(result, extractedValues)
		}
		selectRow = true
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
				"value set has incorrect number of values: %s",
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
			failedIndex := req.Values[i]
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
