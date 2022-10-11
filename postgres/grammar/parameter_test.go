package grammar

import (
	"reflect"
	"testing"
	"time"
)

type testColumn struct {
	field string
	value interface{}
}

func (t testColumn) AutoIncrement() bool {
	return false
}

func (t testColumn) PrimaryKey() bool {
	return false
}

func (t testColumn) Field() string {
	return t.field
}

func (t testColumn) Name() string {
	return t.field
}

func (t testColumn) Value(data interface{}) interface{} {
	return t.value
}

func (t testColumn) Type() reflect.Type {
	return reflect.TypeOf(t.value)
}

type testPerson struct{}

func Test_Parameter_Replace(test *testing.T) {
	test.Run("Success", func(t *testing.T) {
		p := parameter{
			column: &testColumn{
				field: "name",
				value: "ctl",
			},
		}
		sql := "@name"
		sql = p.Replace(sql, testPerson{})
		if sql != "'ctl'" {
			t.Fatal("err: " + sql)
		}
	})
}

func Test_Time(t *testing.T) {
	currentTime := time.Now()
	t.Fatal(reflect.TypeOf(currentTime).Name())

}
