package metadata

import (
	"reflect"
	"strings"
	"sync"

	"github.com/xm-chentl/goresource"
)

const (
	TagName          = "postgres" // 字段名
	TagPrimary       = "pk"       // 主键
	TagAutoIncrement = "auto"     // 自增
)

var (
	rw          sync.RWMutex
	nameOfTable = make(map[string]ITable)
)

type table struct {
	columns             []IColumn
	entryType           reflect.Type
	tableName           string
	primaryKeyColumn    IColumn
	autoIncrementColumn IColumn
}

func (t table) Name() string {
	return t.tableName
}

func (t table) Columns() []IColumn {
	return t.columns
}

func (t table) ColumnMap() (res map[string]IColumn) {
	res = make(map[string]IColumn)
	for _, c := range t.columns {
		res[c.Field()] = c
	}

	return
}

func (t table) GetColumn(key string) (res IColumn) {
	for _, v := range t.columns {
		if strings.EqualFold(key, v.Field()) {
			res = v
			break
		}
	}

	return
}

func (t *table) PrimaryKeyColumn() IColumn {
	if t.primaryKeyColumn == nil {
		// todo: 如果存在不设置的，每次都会遍历
		for index, c := range t.columns {
			if c.PrimaryKey() {
				t.primaryKeyColumn = t.columns[index]
				break
			}
		}
	}

	return t.primaryKeyColumn
}

func (t *table) AutoIncrementColumn() IColumn {
	if t.autoIncrementColumn == nil {
		// todo: 如果存在不设置的，每次都会遍历
		for index, c := range t.columns {
			if c.AutoIncrement() {
				t.autoIncrementColumn = t.columns[index]
				break
			}
		}
	}

	return t.autoIncrementColumn
}

func Get(entry goresource.IDbModel) ITable {
	rw.RLock()
	inst, ok := nameOfTable[entry.Table()]
	if ok {
		defer rw.RUnlock()
		return inst
	}
	rw.RUnlock()

	rw.Lock()
	defer rw.Unlock()
	entryRt := reflect.TypeOf(entry)
	if entryRt.Kind() == reflect.Ptr {
		entryRt = entryRt.Elem()
	}

	tableInst := &table{
		entryType: entryRt,
		columns:   make([]IColumn, 0),
		tableName: entry.Table(),
	}
	recursionNestedStruct(entryRt, &tableInst.columns)
	nameOfTable[entry.Table()] = tableInst

	return tableInst
}

// recursionNestedStruct 递归嵌套结构
func recursionNestedStruct(structType reflect.Type, columns *[]IColumn) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		_, ok := field.Tag.Lookup(TagName)
		if field.Type.Kind() == reflect.Struct && !ok {
			recursionNestedStruct(field.Type, columns)
		} else if ok {
			*columns = append(*columns, &column{
				field: field,
			})
		}
	}
}
