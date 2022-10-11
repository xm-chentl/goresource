package metadata

import (
	"reflect"
	"testing"
)

func Test_Value(test *testing.T) {
	test.Run("Success", func(t *testing.T) {
		entry := testPerson{
			ID:   "test-001",
			Name: "test-name",
		}
		entryRt := reflect.TypeOf(entry)
		nameColumn := column{}
		nameColumn.field, _ = entryRt.FieldByName("Name")
		nameValue := nameColumn.Value(entry)
		if nameValue == nil || nameValue.(string) != entry.Name {
			t.Fatal("err")
		}
	})
}
