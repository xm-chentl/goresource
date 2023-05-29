package mysqlex

import (
	"fmt"

	"gorm.io/gorm"
)

type IOption interface {
	Apply(*gorm.DB) *gorm.DB
}

type OptionQuerySelect struct {
	Fields []string
}

func (s OptionQuerySelect) Apply(db *gorm.DB) *gorm.DB {
	return db.Select(s.Fields)
}

func NewOptionQuerySelect(fields ...string) IOption {
	return &OptionQuerySelect{
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
