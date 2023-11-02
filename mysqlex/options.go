package mysqlex

import (
	"fmt"

	"gorm.io/gorm"
)

type IOption interface {
	Apply(*gorm.DB) *gorm.DB
}

type OptionGroupBy struct {
	Fields []string
}

func (g OptionGroupBy) Apply(db *gorm.DB) *gorm.DB {
	if len(g.Fields) > 0 {
		for _, f := range g.Fields {
			db = db.Group(f)
		}
	}

	return db
}

func NewOptionGroupBy(fields ...string) IOption {
	return &OptionGroupBy{
		Fields: fields,
	}
}

type OptionSelectField struct {
	Fields []string
}

func (s OptionSelectField) Apply(db *gorm.DB) *gorm.DB {
	return db.Select(s.Fields)
}

func NewOptionSelectField(fields ...string) IOption {
	return &OptionSelectField{
		Fields: fields,
	}
}

type OptionSaveOmit struct {
	Fields []string
}

func (s OptionSaveOmit) Apply(db *gorm.DB) *gorm.DB {
	return db.Omit(s.Fields...)
}

func NewOptionSaveOmit(fields ...string) IOption {
	return &OptionSaveOmit{
		Fields: fields,
	}
}

type OptionTableSuffix struct {
	Value string
}

func (t OptionTableSuffix) Apply(db *gorm.DB) *gorm.DB {
	return db.Table(t.Value)
}

func NewOptionTableSuffix(tableName, format string, args ...interface{}) IOption {
	return &OptionTableSuffix{
		Value: fmt.Sprintf(tableName+format, args...),
	}
}
