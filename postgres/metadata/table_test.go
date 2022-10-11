package metadata

import "testing"

type testPerson struct {
	ID   string `postgres:"id"`
	Name string `postgres:"name"`
}

func (t testPerson) Table() string {
	return "person"
}

func (t testPerson) GetID() interface{} {
	return t.ID
}

func (t *testPerson) SetID(v interface{}) {
	if idv, ok := v.(string); ok {
		t.ID = idv
	}
}

func Test_Get(test *testing.T) {
	test.Run("table.name", func(t *testing.T) {
		person := testPerson{}
		tableInst := Get(&person)
		if tableInst.Name() != "person" {
			t.Fatal("err")
		}
	})

	test.Run("table.columns", func(t *testing.T) {
		person := testPerson{}
		tableInst := Get(&person)
		if len(tableInst.Columns()) != 2 {
			t.Fatal("err")
		}
		if tableInst.Columns()[0].Field() != "id" || tableInst.Columns()[1].Field() != "name" {
			t.Fatal("err")
		}
	})
}
