package column

import (
	"fmt"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table"
	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

type column struct {
	Name     *t.Token
	Elements []*t.Token
}

type ColumnTable struct {
	Columns      map[int]*column
	Name         *t.Token
	columnTypes  []t.TokenKind
	orderByIndex int
}

func (ct ColumnTable) IsInitialized() bool {
	return len(ct.Columns) > 0
}

func validateRequest(req *parser.CreateTableQuery) error {
	// Check if all columns has different names
	for i := 0; i < len(*req.Cols)-1; i++ {
		for j := i + 1; j < len(*req.Cols); j++ {
			if ((*req.Cols)[i]).Equals((*req.Cols)[j]) {
				return fmt.Errorf("columns are not unique")
			}
		}
	}
	// Check if order by column name in column names
	var orderValid bool = false
	for i := 0; i < len(*req.Cols); i++ {
		if (*req.Cols)[i].Name.Value == req.OrderBy.Value {
			orderValid = true
		}
	}
	if !orderValid {
		return fmt.Errorf("no column specified in ORDER BY was found")
	}
	return nil
}

func (ct ColumnTable) Create(req *parser.CreateTableQuery) (table.Table, string, error) {
	err := validateRequest(req)
	if err != nil {
		return nil, "", err
	}
	typeCorreation := map[string]t.TokenKind{
		"int":  t.NumericKind,
		"text": t.IdentifierKind,
	}

	ct.Name = req.Name

	ct.Columns = make(map[int]*column)
	for i, columnDefinition := range *req.Cols {
		tempColumn := &column{
			Name:     columnDefinition.Name,
			Elements: []*t.Token{},
		}
		ct.columnTypes = append(ct.columnTypes, typeCorreation[columnDefinition.Datatype.Value])
		if columnDefinition.Name.Value == req.OrderBy.Value {
			ct.orderByIndex = i
		}
		ct.Columns[i] = tempColumn
	}
	return table.Table(ct), ct.Name.Value, nil
}

func (ct ColumnTable) Select(req *parser.SelectQuery) ([][]*t.Token, error) {
	// NOTE: Only selects without WHERE allowed now
	// tableColumns := ct.GetColumnsNames()
	indexCorrelation := make(map[string]int)
	index := 0
	for _, column := range ct.Columns {
		indexCorrelation[column.Name.Value] = index
		index++
	}

	indexesToSelect := []int{}
	for _, columnToSelect := range req.Columns {
		if index, ok := indexCorrelation[columnToSelect.Value]; !ok {
			return nil, fmt.Errorf("column with name %s is not represented inside table (actual columns: [%s])",
				columnToSelect,
				ct.GetColumns(),
			)
		} else {
			indexesToSelect = append(indexesToSelect, index)
		}
	}

	header := make([]*t.Token, len(indexesToSelect))
	index = 0
	for i, extractedHeaderPosition := range indexesToSelect {
		header[i] = ct.Columns[extractedHeaderPosition].Name
	}
	result := [][]*t.Token{header}
	// Index is a number of row
	// Element is an element of row under ORDER BY column
	for index := range ct.Columns[ct.orderByIndex].Elements {
		// here must be condition on order by element
		var selectRow bool = true
		extractedValues := []*t.Token{}
		if selectRow {
			// get element with index from all columns
			for _, list := range ct.Columns {
				extractedValues = append(extractedValues, list.Elements[index])
			}
			result = append(result, extractedValues)
		}
	}
	return result, nil
}

func (ct ColumnTable) Insert(req *parser.InsertIntoQuery) (table.Table, error) {
	for _, set := range req.Values {
		for valueIndex := range set {
			if ct.columnTypes[valueIndex] != set[valueIndex].Kind {
				return nil, fmt.Errorf("token %s has incorrect type (expected: %s)", set[valueIndex].Value, t.KindToString(ct.columnTypes[valueIndex]))
			}
		}
		// Find position to insert
		position, err := getInsertionPosition(set[ct.orderByIndex], ct.Columns[ct.orderByIndex].Elements)
		if err != nil {
			return nil, fmt.Errorf("an error occured while searching pos searching for value %s in %s column", set[ct.orderByIndex].Value, ct.Columns[ct.orderByIndex].Name.Value)
		}
		// Insert values
		colIndex := 0
		for _, column := range ct.Columns {
			if position == 0 {
				column.Elements = append([]*t.Token{set[colIndex]}, column.Elements...)
			} else if position == len(column.Elements) {
				column.Elements = append(column.Elements, set[colIndex])
			} else {
				column.Elements = append(column.Elements[:position+1], column.Elements[position:]...)
				column.Elements[position] = set[colIndex]
			}
			colIndex++
		}
	}

	return table.Table(ct), nil
}

func (ct ColumnTable) Engine() string {
	return "Column"
}

func (ct ColumnTable) ShowCreate() string {
	return ""
}

func (ct ColumnTable) GetColumns() string {
	var result string
	if !ct.IsInitialized() {
		return ""
	}
	index := 0
	for _, column := range ct.Columns {
		result += fmt.Sprintf("(%s, %s), ", column.Name.Value, t.KindToString(ct.columnTypes[index]))
		index++
	}
	return result[:len(result)-2]
}

func (ct ColumnTable) Count() int {
	return len(ct.Columns[ct.orderByIndex].Elements)
}

func (ct ColumnTable) GetColumnsNames() []string {
	result := make([]string, len(ct.Columns))
	index := 0
	for _, column := range ct.Columns {
		result[index] = column.Name.Value
		index++
	}
	return result
}

func getInsertionPosition(compare *t.Token, values []*t.Token) (int, error) {
	start, end := 0, len(values)
	for start < end {
		mid := (start + end) / 2
		var (
			val int
			err error
		)
		val, err = values[mid].Compare(compare)
		if err != nil {
			return -2, err
		}
		if val == 0 {
			return mid, nil
		}
		val, err = values[mid].Compare(compare)
		if err != nil {
			return -2, err
		}
		if val == -1 {
			start = mid + 1
		}
		val, err = values[mid].Compare(compare)
		if err != nil {
			return -2, err
		}
		if val == 1 {
			end = mid
		}
	}
	return end, nil
}
