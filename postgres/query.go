package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/errs"
	"github.com/xm-chentl/goresource/postgres/grammar"
	"github.com/xm-chentl/goresource/postgres/metadata"

	"github.com/jackc/pgtype"
)

type query struct {
	ctx       context.Context
	pool      *pool
	fields    []string
	where     string
	whereArgs []interface{}
	page      int
	pageSize  int
	orders    []string
	orderBys  []string
	opts      []interface{}
}

func (q *query) Count(entry goresource.IDbModel) (res int64, err error) {
	defer q.reset()

	sql, args := grammar.Count(metadata.Get(entry), q.getArgs()...)
	conn, err := q.pool.getConn()
	if err != nil {
		return
	}
	defer conn.Release()

	row := conn.QueryRow(q.ctx, sql, args...)
	err = row.Scan(&res)

	return
}

func (q query) Exec(res interface{}, args ...interface{}) (err error) {
	defer q.reset()

	if len(args) == 0 {
		err = errs.QueryArgsError
		return
	}

	resRt := reflect.TypeOf(res)
	resRv := reflect.ValueOf(res)
	if resRt.Kind() == reflect.Ptr {
		if resRt.Elem().Kind() == reflect.Slice {
			resRt = resRt.Elem().Elem()
		} else {
			err = errs.ResIsNotSlice
			return
		}
	} else {
		err = errs.ResIsNotPtr
		return
	}
	err = q.scan(resRt, resRv, args[0].(string), args[1:]...)

	return
}

func (q *query) Fields(fields ...interface{}) goresource.IQuery {
	if len(fields) > 0 {
		for _, field := range fields {
			if v, ok := field.(string); ok {
				if strings.Contains(v, `"`) {
					q.fields = append(q.fields, v)
					continue
				}
				q.fields = append(q.fields, metadata.FormatField(v))
			}
		}
	}
	return q
}

func (q query) First(res interface{}) (err error) {
	resRt := reflect.TypeOf(res)
	if resRt.Kind() != reflect.Ptr {
		err = errs.ResIsNotPtr
		return
	} else {
		resRt = resRt.Elem()
	}
	if resRt.Kind() != reflect.Struct {
		err = errs.ResIsNotStruct
		return
	}

	resRv := reflect.ValueOf(res)
	resRvSlice := reflect.New(
		reflect.SliceOf(resRt),
	)
	if err := q.queryData(resRt, resRvSlice); err != nil {
		return err
	}
	if resRvSlice.Elem().Len() > 0 {
		resRv.Elem().Set(resRvSlice.Elem().Index(0))
	}

	return
}

func (q query) Find(res interface{}) (err error) {
	resRt := reflect.TypeOf(res)
	resRv := reflect.ValueOf(res)
	if resRt.Kind() == reflect.Ptr {
		if resRt.Elem().Kind() == reflect.Slice {
			resRt = resRt.Elem().Elem()
		} else {
			err = errs.ResIsNotSlice
			return
		}
	} else {
		err = errs.ResIsNotPtr
		return
	}
	if err = q.queryData(resRt, resRv); err != nil {
		return
	}

	return
}

func (q query) ToArray(res interface{}) error {
	return q.Find(res)
}

func (q *query) Where(args ...interface{}) goresource.IQuery {
	if len(args) > 0 {
		q.where = args[0].(string)
	}
	if len(args) > 1 {
		q.whereArgs = args[1:]
	}

	return q
}

func (q *query) Page(page int) goresource.IQuery {
	q.page = page
	if q.page < 1 {
		q.page = 1
	}

	return q
}

func (q *query) PageSize(pageSize int) goresource.IQuery {
	q.pageSize = pageSize
	if q.pageSize < 1 {
		q.pageSize = 20
	}

	return q
}

func (q *query) Asc(fields ...string) goresource.IQuery {
	for _, field := range fields {
		q.orders = append(q.orders, fmt.Sprintf(`"%s"`, field))
	}

	return q
}

func (q *query) Desc(fields ...string) goresource.IQuery {
	for _, field := range fields {
		q.orderBys = append(q.orderBys, fmt.Sprintf(`"%s"`, field))
	}

	return q
}

func (q *query) SetOpts(opts ...interface{}) goresource.IQuery {
	for _, o := range opts {
		if o != nil {
			q.opts = opts
		}
	}

	return q
}

func (q query) getArgs() (args []interface{}) {
	args = make([]interface{}, 0)
	if strings.TrimSpace(q.where) == "" {
		return
	}
	args = append(args, q.where)
	args = append(args, q.whereArgs...)

	return
}

func (q *query) queryData(rt reflect.Type, resultsOfRv reflect.Value) (err error) {
	defer q.reset()

	sql, args := grammar.Select(
		metadata.Get(
			reflect.New(rt).Interface().(goresource.IDbModel),
		),
		q.fields,
		q.getArgs()...,
	)
	// Todo: 后续封装至grammar
	if len(q.orders) > 0 {
		sql += fmt.Sprintf(" ORDER BY %s ASC", strings.Join(q.orders, ", "))
	}
	if len(q.orderBys) > 0 {
		if len(q.orders) == 0 {
			sql += " ORDER BY"
		} else {
			sql += ", "
		}
		sql += fmt.Sprintf(" %s DESC", strings.Join(q.orderBys, ", "))
	}
	if q.page > 0 || q.pageSize > 0 {
		if q.page == 0 {
			q.page = 1
		}
		if q.pageSize == 0 {
			q.pageSize = 1
		}
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", q.pageSize, ((q.page - 1) * q.pageSize))
	}

	err = q.scan(rt, resultsOfRv, sql, args...)

	return
}

func (q query) scan(rt reflect.Type, resultsOfRv reflect.Value, sql string, args ...interface{}) (err error) {
	conn, err := q.pool.getConn()
	if err != nil {
		return
	}
	defer conn.Release()
	rows, err := conn.Query(q.ctx, sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	results := reflect.MakeSlice(reflect.SliceOf(rt), 0, 0)
	rv := reflect.New(rt).Elem()
	bindFieldMap := make(map[string]interface{})
	nestedBindStructByMap(rt, rv, bindFieldMap)
	fieldDescArray := rows.FieldDescriptions()
	var resArray []interface{}
	for rows.Next() {
		resArray, err = rows.Values()
		if err != nil {
			return
		}
		for index := range fieldDescArray {
			fieldDes := fieldDescArray[index]
			key := strings.ToLower(string(fieldDes.Name))
			_, ok := bindFieldMap[key]
			if !ok {
				continue
			}

			bResRv := reflect.ValueOf(bindFieldMap[key])
			resRv := reflect.ValueOf(resArray[index])
			if resRv.IsValid() {
				resRvValue := reflect.New(resRv.Type())
				resRvValue.Elem().Set(resRv)
				if pgValue, ok := resRvValue.Interface().(pgtype.Value); ok {
					pgValue.AssignTo(bResRv.Interface())
				} else {
					bResRv.Elem().Set(resRv)
				}
			}
		}
		results = reflect.Append(results, rv)
	}
	resultsOfRv.Elem().Set(results)

	return
}

func (q *query) reset() {
	q.where = ""
	q.whereArgs = []interface{}{}
	q.page = 0
	q.pageSize = 0
	q.orderBys = make([]string, 0)
	q.orders = make([]string, 0)
}

func (q query) exec(res interface{}, args ...interface{}) (err error) {
	defer q.reset()

	if len(args) == 0 {
		err = fmt.Errorf("args parameter error")
		return
	}

	resRt := reflect.TypeOf(res)
	resRv := reflect.ValueOf(res)
	if resRt.Kind() == reflect.Ptr {
		if resRt.Elem().Kind() == reflect.Slice {
			resRt = resRt.Elem().Elem()
		} else {
			err = errors.New("res is not slice")
			return
		}
	} else {
		err = errors.New("res is not ptr")
		return
	}

	sql := args[0].(string)
	args = args[1:]
	conn, err := q.pool.getConn()
	if err != nil {
		return
	}
	defer conn.Release()

	rows, err := conn.Query(q.ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rv := reflect.New(resRt).Elem()
	bindRvs := make([]interface{}, 0)
	nestedBindStruct(resRt, rv, &bindRvs)
	results := reflect.MakeSlice(reflect.SliceOf(resRt), 0, 0)
	var resArray []interface{}
	for rows.Next() {
		resArray, err = rows.Values()
		if err != nil {
			return
		}
		for index := range bindRvs {
			if index+1 > len(resArray) {
				continue
			}

			bResRv := reflect.ValueOf(bindRvs[index])
			resRv := reflect.ValueOf(resArray[index])
			if resRv.IsValid() {
				resRvValue := reflect.New(resRv.Type())
				resRvValue.Elem().Set(resRv)
				if pgValue, ok := resRvValue.Interface().(pgtype.Value); ok {
					pgValue.AssignTo(bResRv.Interface())
				} else {
					bResRv.Elem().Set(resRv)
				}
			}
		}
		results = reflect.Append(results, rv)
	}
	resRv.Elem().Set(results)

	return
}

func nestedBindStruct(resRt reflect.Type, resRv reflect.Value, bindRvs *[]interface{}) {
	for i := 0; i < resRt.NumField(); i++ {
		field := resRt.Field(i)
		_, ok := field.Tag.Lookup(metadata.TagName)
		if field.Type.Kind() == reflect.Struct && !ok && !strings.EqualFold(field.Type.Name(), "time") {
			resRv.Field(i).Set(reflect.New(field.Type).Elem())
			nestedBindStruct(field.Type, resRv.Field(i), bindRvs)
		} else if ok {
			*bindRvs = append(*bindRvs, resRv.FieldByName(field.Name).Addr().Interface())
		}
	}
}

func nestedBindStructByMap(resRt reflect.Type, resRv reflect.Value, bindFields map[string]interface{}) {
	for i := 0; i < resRt.NumField(); i++ {
		field := resRt.Field(i)
		tagName, ok := field.Tag.Lookup(metadata.TagName)
		// 嵌套结构 排队时间字段
		if field.Type.Kind() == reflect.Struct && !ok && !strings.EqualFold(field.Type.Name(), "time") {
			resRv.Field(i).Set(reflect.New(field.Type).Elem())
			nestedBindStructByMap(field.Type, resRv.Field(i), bindFields)
		} else if ok {
			bindFields[strings.ToLower(tagName)] = resRv.FieldByName(field.Name).Addr().Interface()
		} else if !ok {
			bindFields[strings.ToLower(field.Name)] = resRv.FieldByName(field.Name).Addr().Interface()
		}
	}
}
