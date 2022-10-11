package grammar

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/xm-chentl/goresource/postgres/metadata"
)

type IParameter interface {
	Replace(sql string, entry interface{}) string
}

type parameter struct {
	column metadata.IColumn
}

func (d parameter) Replace(sql string, entry interface{}) string {
	var replaceValue string
	switch d.column.Type().Kind() {
	case reflect.String, reflect.Struct:
		replaceValue = fmt.Sprintf("'%v'", d.column.Value(entry))
	default:
		replaceValue = fmt.Sprintf("%v", d.column.Value(entry))
	}

	return strings.ReplaceAll(sql, VarTag+d.column.Field(), replaceValue)
}

func NewParameter(column metadata.IColumn) IParameter {
	return &parameter{
		column: column,
	}
}
