package parser

import (
	"encoding/json"
	"fmt"
	"strconv"

	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// SelectQuery give access to requests with syntax
//
// SELECT (<column_name1>[...]) FROM <table_name>;
//
// where <column_nameX> and <table_name> must be tokens with IdentifierKind
type SelectQuery struct {
	Columns    []*t.Token   `json:"item"`
	From       *t.Token     `json:"from"`
	Conditions []*Condition `json:"where"`
}

// String method needs to be implemented in order to implement Query interface.
// Returns JSON object describing necessary information
func (slct SelectQuery) String() string {
	bytes, _ := json.Marshal(slct)
	return string(bytes)
}

// Equals method needs to be implemented in order to implement Query interface.
// Returns true if tokens for columns and table names are equal.
func (slct SelectQuery) Equals(other *SelectQuery) bool {
	if len(slct.Columns) != len(other.Columns) {
		return false
	}
	if len(slct.Conditions) != len(other.Conditions) {
		return false
	}
	for index := range slct.Columns {
		if !slct.Columns[index].Equals(other.Columns[index]) {
			return false
		}
	}
	for index := range slct.Conditions {
		if !slct.Conditions[index].Equals(other.Conditions[index]) {
			return false
		}
	}
	return slct.From.Equals(other.From)
}

// CreateOriginal method needs to be implemented in order to implement Query interface.
// Returns original SQL query representing data in current Query
func (slct SelectQuery) CreateOriginal() string {
	result := fmt.Sprintf("SELECT %s FROM %s;",
		t.Bracketize(slct.Columns),
		slct.From.Value,
	)
	return result
}

type Condition struct {
	Column          *t.Token `json:"column_name"`
	Value           *t.Token `json:"value"`
	ConditionSymbol *t.Token `json:"condition"`
}

func (c *Condition) Equals(other *Condition) bool {
	if !c.Column.Equals(other.Column) {
		return false
	}
	if !c.Value.Equals(other.Value) {
		return false
	}
	if !c.ConditionSymbol.Equals(other.ConditionSymbol) {
		return false
	}
	return true
}

func (c *Condition) EvaluateWithValues(left *t.Token) (bool, error) {
	var (
		lval             int
		rval             int
		stringProcessors = map[string]func(string, string) bool{
			"==": func(s1, s2 string) bool { return s1 == s2 },
			">":  func(s1, s2 string) bool { return s1 > s2 },
			"<":  func(s1, s2 string) bool { return s1 < s2 },
		}
		numericProcessors = map[string]func(int, int) bool{
			"==": func(s1, s2 int) bool { return s1 == s2 },
			">":  func(s1, s2 int) bool { return s1 > s2 },
			"<":  func(s1, s2 int) bool { return s1 < s2 },
		}
	)
	if c.Value.Kind != left.Kind {
		return false, fmt.Errorf("cannot compare tokens with different kinds (%s, %s => %s, %s", c.Value, left.Value, t.KindToString(c.Value.Kind), t.KindToString(left.Kind))
	}
	if c.Value.Kind == t.NumericKind {
		rval, _ = strconv.Atoi(c.Value.Value)
		lval, _ = strconv.Atoi(left.Value)
		return numericProcessors[c.ConditionSymbol.Value](lval, rval), nil
	} else {
		return stringProcessors[c.ConditionSymbol.Value](left.Value, c.Value.Value), nil
	}
}
