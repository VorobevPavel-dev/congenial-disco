package table

import (
	"testing"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
)

func TestLinearTable(t *testing.T) {
	t.Run("Test linear table creation", func(t *testing.T) {
		request := "create table test (id int, name text);"
		statement := parser.Parse(request)
		if statement == nil {
			t.Error("Unable to parse input string as any kind of requests")
		}
		if statement.CreateTableStatement == nil {
			t.Error("Unable to parse iwnput string as CREATE TABLE request")
		}
		table := &LinearTable{}
		_, _ = table.Create(statement.CreateTableStatement)
		if !table.IsInitialized() {
			t.Error("Table has no columns inside after creation request")
		}
		// recreatedRequest := table.ShowCreate()
		// if request != recreatedRequest {
		// 	t.Errorf("Error recreating request. Expected: %s, got: %s", request, recreatedRequest)
		// }
	})
}
