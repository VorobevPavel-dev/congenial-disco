package session

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	token "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

type command struct {
	SQL            string `json:"sql"`
	ExpectedOutput string `json:"expected_output"`
	ExpectError    bool   `json:"expect_error"`
}

func TestCommandSequenceExecution(t *testing.T) {
	var (
		session           Session = InitSession()
		inputSequenceFile string  = "./input_data.json"
		rawSequence       []byte
		inputSequence     []command
		index             int
		command           command
	)
	defer func(pos *int) {
		if err := recover(); err != nil {
			t.Errorf("Unexpected panic on set %d: %v", index+1, err)
		}
	}(&index)
	rawSequence, _ = ioutil.ReadFile(inputSequenceFile)
	json.Unmarshal(rawSequence, &inputSequence)
	for index, command = range inputSequence {
		result, err := session.ExecuteCommand(command.SQL)
		if err != nil {
			if command.ExpectError {
				continue
			}
			t.Errorf("error on command #%d: %v", index+1, err)
		}
		diff := strings.Compare(strings.TrimSpace(result), command.ExpectedOutput)
		if diff != 0 {
			t.Errorf(
				"actial result differs from expected on set %d: %s => %s",
				index+1,
				strings.TrimSpace(result),
				command.ExpectedOutput,
			)
		}
	}
}

func BenchmarkMassiveTableCreation(b *testing.B) {
	session := InitSession()
	inputs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		// Generating random CREATE TABLE query
		tableName := token.GenerateRandomToken(token.IdentifierKind)
		//Generate column definition
		columns := make([]parser.ColumnDefinition, 10)
		for i := range columns {
			columns[i] = parser.ColumnDefinition{
				Name:     token.GenerateRandomToken(token.IdentifierKind),
				Datatype: token.GenerateRandomToken(token.TypeKind),
			}
		}
		inputQuery := &parser.CreateTableQuery{
			Name: tableName,
			Cols: &columns,
		}
		inputs[i] = (*inputQuery).CreateOriginal()
	}
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := session.ExecuteCommand(inputs[i])
		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
	}
}
