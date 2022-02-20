package linear

import (
	"testing"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table"
)

func TestLinearTable(t *testing.T) {
	t.Run("Test linear table creation", func(t *testing.T) {
		request := "create table test (id int, name text) engine linear;"
		statement, err := parser.Parse(request)
		if err != nil {
			t.Error(err)
		}
		if statement == nil {
			t.Error("Unable to parse input string as any kind of requests")
			return
		} else if statement.CreateTableStatement == nil {
			t.Error("Unable to parse input string as CREATE TABLE request")
		}
		table := table.Table(&LinearTable{})
		table, _, _ = table.Create(statement.CreateTableStatement)
		if !table.IsInitialized() {
			t.Error("Table has no columns inside after creation request")
		}
		// recreatedRequest := table.ShowCreate()
		// if request != recreatedRequest {
		// 	t.Errorf("Error recreating request. Expected: %s, got: %s", request, recreatedRequest)
		// }
	})
}
