package grammar

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/postgres/metadata"
)

// Select 生成查询语句 fields 指定查询字段 args 0 where > 1 where-args
func Select(table metadata.ITable, fields []string, args ...interface{}) (sql string, newArgs []interface{}) {
	var bf bytes.Buffer
	bf.WriteString("SELECT ")
	fieldArray := make([]string, 0)
	if len(fields) > 0 {
		columnMap := table.ColumnMap()
		for _, f := range fields {
			if _, ok := columnMap[f]; ok {
				fieldArray = append(fieldArray, f)
			}
		}
	} else {
		for _, c := range table.Columns() {
			fieldArray = append(fieldArray, c.Field())
		}
	}

	bf.WriteString(strings.Join(fieldArray, ", "))
	bf.WriteString(" FROM ")
	bf.WriteString(table.Name())
	if len(args) > 0 {
		where, whereArges := Where(args...)
		bf.WriteString(where)
		newArgs = whereArges
	}
	sql = bf.String()

	return
}

// Insert 生成插入语句
func Insert(table metadata.ITable, entry goresource.IDbModel) (sql string, args []interface{}) {
	var bf bytes.Buffer
	bf.WriteString("INSERT INTO ")
	bf.WriteString(table.Name())
	columnArray := make([]string, 0)
	varArray := make([]string, 0)
	args = make([]interface{}, 0)
	index := 1
	for _, column := range table.Columns() {
		if column.AutoIncrement() {
			continue
		}

		columnArray = append(columnArray, column.Field())
		varArray = append(varArray, fmt.Sprintf("$%d", index))
		args = append(args, column.Value(entry))
		index++
	}
	bf.WriteString(" (")
	bf.WriteString(strings.Join(columnArray, ", "))
	bf.WriteString(") VALUES (")
	bf.WriteString(strings.Join(varArray, ", "))
	bf.WriteString(");")
	sql = bf.String()
	return
}

// Update 生成更新语句 args 0 where > 1 where-args
func Update(table metadata.ITable, entry goresource.IDbModel, fields []string, args ...interface{}) (sql string, newArgs []interface{}) {
	var bf bytes.Buffer
	bf.WriteString("UPDATE ")
	bf.WriteString(table.Name())
	newArgs = make([]interface{}, 0)
	updateFields := make([]string, 0)
	columns := make([]metadata.IColumn, 0)
	if len(fields) == 0 {
		for index, column := range table.Columns() {
			if column.AutoIncrement() || column.PrimaryKey() {
				continue
			}
			columns = append(columns, table.Columns()[index])
		}
	} else {
		for _, field := range fields {
			column, ok := table.ColumnMap()[fmt.Sprintf(`"%s"`, field)]
			// 过滤 不存在，自增的
			if !ok || column.AutoIncrement() || column.PrimaryKey() {
				continue
			}
			columns = append(columns, column)
		}
	}
	for index := range columns {
		c := columns[index]
		updateFields = append(updateFields, fmt.Sprintf(`%s=$%d`, c.Field(), index+1))
		newArgs = append(newArgs, c.Value(entry))
	}
	bf.WriteString(" SET ")
	bf.WriteString(strings.Join(updateFields, ","))
	if len(args) > 0 {
		where, whereArgs := Where(args...)
		bf.WriteString(where)
		if len(whereArgs) > 0 {
			newArgs = append(newArgs, whereArgs...)
		}
	}
	bf.WriteString(";")
	sql = bf.String()

	return
}

// Delete 生成删除语句 args 0 where > 1 where-args
func Delete(table metadata.ITable, entry goresource.IDbModel, args ...interface{}) (sql string, newArgs []interface{}) {
	var bf bytes.Buffer
	bf.WriteString("DELETE FROM ")
	bf.WriteString(table.Name() + " ")
	if len(args) > 0 {
		where, whereArgs := Where(args...)
		bf.WriteString(where)
		newArgs = make([]interface{}, 0)
		newArgs = append(newArgs, whereArgs...)
	}
	bf.WriteString(";")
	sql = bf.String()

	return
}

// Count 生成count语句 args 0 where > 1 where-args
func Count(table metadata.ITable, args ...interface{}) (sql string, newArgs []interface{}) {
	var bf bytes.Buffer
	bf.WriteString("SELECT count(1) FROM ")
	bf.WriteString(table.Name())
	newArgs = make([]interface{}, 0)
	if len(args) > 0 {
		where, whereArgs := Where(args...)
		bf.WriteString(where)
		newArgs = append(newArgs, whereArgs...)
	}
	sql = bf.String()

	return
}

// Where 生成条件语句 args 0 where > 1 where-args
func Where(args ...interface{}) (sql string, newArgs []interface{}) {
	if len(args) == 0 {
		return
	}

	where := args[0].(string)
	sql = fmt.Sprintf(" WHERE %s", where)
	if len(args) > 1 {
		newArgs = args[1:]
	}

	return
}
