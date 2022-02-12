package table

import (
	"fmt"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
	"github.com/VorobevPavel-dev/congenial-disco/utility"
)

type LinearTable struct {
	Columns  []*parser.ColumnDefinition
	Elements [][]*tokenizer.Token
	Name     *tokenizer.Token
}

func (lt LinearTable) IsInitialized() bool {
	return len(lt.Columns) > 0
}

// Create table will initialize columns inside fresh linear table
// It will take parser.ColumnDefinitions from request and append them to LinearTable.Columns slice
// Returns table name if all happened without errors
func (lt LinearTable) Create(req *parser.CreateTableQuery) (Table, string, error) {
	lt.Columns = make([]*parser.ColumnDefinition, len(req.Cols))
	for i := range req.Cols {
		lt.Columns[i] = req.Cols[i]
	}
	lt.Name = req.Name
	return Table(lt), lt.Name.Value, nil
}

func (lt LinearTable) Select(req *parser.SelectStatement) (*[][]tokenizer.Token, error) {
	return nil, nil
}

func (lt LinearTable) Insert(req *parser.InsertIntoQuery) (Table, error) {
	// typeCorreation[column.Datatype.Value] will give Kind which is supported on current column
	typeCorreation := map[string]tokenizer.TokenKind{
		"int":  tokenizer.NumericKind,
		"text": tokenizer.IdentifierKind,
	}
	// Build map with <column name> - <supported token type>
	// For query CREATE TABLE test (id INT, name TEXT);
	// columnCorrelation will have next elements:
	//		"id": tokenizer.NumericKind
	//		"name": tokenizer.IdentifierKind
	columnCorrelation := make(map[string]tokenizer.TokenKind)
	for _, column := range lt.Columns {
		columnCorrelation[column.Name.Value] = typeCorreation[column.Datatype.Value]
	}

	mapToInsert := make(map[string]*tokenizer.Token)
	for _, column := range lt.Columns {
		mapToInsert[column.Name.Value] = &tokenizer.Token{
			Value: "null",
			Kind:  tokenizer.IdentifierKind,
		}
	}

	// If names of columns were not provided - whole reow must be inserted
	if len(req.ColumnNames) != 0 {
		// Assert that all columns in request are actually in table and in correct order
		lastVisited := -1
		columnNames := lt.GetColumnsNames()

		for index, reqColumn := range req.ColumnNames {
			tempVisited := utility.FindStringInSlice(columnNames, reqColumn.Value)
			if tempVisited == -1 {
				return nil, fmt.Errorf("column name %s was not found in table %s",
					reqColumn.Value,
					req.Table.Value,
				)
			}
			if tempVisited < lastVisited {
				return nil, fmt.Errorf("columns in request %s are in incorrect order (desired: %s)",
					req.CreateOriginal(),
					utility.StringSliceToString(lt.GetColumnsNames()),
				)
			}
			lastVisited = tempVisited
			mapToInsert[reqColumn.Value] = req.Values[index]
		}
	} else {
		// Assert that number of inserting values is equal to number of columns in table
		if len(req.Values) != len(lt.Columns) {
			return nil, fmt.Errorf("unexpected number of values without column specification (expected: %d, got: %d)",
				len(req.Values),
				len(lt.Columns),
			)
		}
		index := 0
		for key := range mapToInsert {
			mapToInsert[key] = req.Values[index]
			index++
		}
	}

	// Check that mapToInsert's values have supported type
	// Note: columnCorrelation has same number of keys as mapToInsert
	for column, supportedType := range columnCorrelation {
		if mapToInsert[column].Kind != supportedType {
			// TODO: Add NULL type for token
			// Skip "null" token
			if mapToInsert[column].Value == "null" {
				continue
			}
			return nil, fmt.Errorf("value %s for column with name %s has unsupported type (expected type to insert: %s)",
				mapToInsert[column].Value,
				column,
				tokenizer.KindToString(int(supportedType)),
			)
		}
	}

	// Get list of values from mapToInsert
	toInsert := make([]*tokenizer.Token, len(mapToInsert))
	index := 0
	for _, value := range mapToInsert {
		toInsert[index] = value
		index++
	}
	lt.Elements = append(lt.Elements, toInsert)
	return Table(lt), nil
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
