package metadata

import (
	"reflect"
	"strings"
)

type column struct {
	field reflect.StructField
}

func (c column) AutoIncrement() (res bool) {
	_, res = c.field.Tag.Lookup(TagAutoIncrement)
	return
}

func (c column) PrimaryKey() (res bool) {
	_, res = c.field.Tag.Lookup(TagPrimary)
	return
}

func (c column) Name() string {
	return c.field.Name
}

func (c column) Field() (res string) {
	value, ok := c.field.Tag.Lookup(TagName)
	if !ok {
		value = strings.ToLower(c.field.Name)
	}
	res = FormatField(value)

	return
}

func (c column) Value(data interface{}) (res interface{}) {
	dataRv := reflect.ValueOf(data)
	if dataRv.Kind() == reflect.Ptr {
		dataRv = dataRv.Elem()
	}
	res = dataRv.FieldByName(c.field.Name).Interface()

	return
}

func (c column) Type() reflect.Type {
	return c.field.Type
}
