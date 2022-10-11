package metadata

import "reflect"

type IColumn interface {
	PrimaryKey() bool
	AutoIncrement() bool
	Field() string
	Name() string
	Value(data interface{}) interface{}
	Type() reflect.Type
}
