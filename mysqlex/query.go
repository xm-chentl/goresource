package mysqlex

import (
	"reflect"
	"strings"

	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/errs"

	"gorm.io/gorm"
)

type query struct {
	db        *gorm.DB
	fields    []string
	whereArgs []interface{}
	order     string
	whereSql  string
	page      int
	pageSize  int
	opts      []interface{}
}

func (q *query) Count(entry goresource.IDbModel) (count int64, err error) {
	defer q.reset()

	db := q.db.Model(entry)
	if q.whereSql != "" {
		db = db.Where(q.whereSql, q.whereArgs...)
	}
	if len(q.opts) > 0 {
		for _, o := range q.opts {
			if v, ok := o.(IOption); ok {
				db = v.Apply(db)
			}
		}
	}
	err = db.Count(&count).Error

	return
}

func (q query) Exec(res interface{}, args ...interface{}) (err error) {
	if len(args) == 0 {
		err = errs.QueryGrammarEmptyError
		return
	}

	sql := args[0].(string)
	sqlArgs := make([]interface{}, 0)
	if len(args) > 1 {
		sqlArgs = append(sqlArgs, args[1:]...)
	}

	err = q.db.Raw(sql, sqlArgs...).Scan(res).Error

	return
}

func (q *query) Fields(args ...interface{}) goresource.IQuery {
	for _, v := range args {
		q.fields = append(q.fields, v.(string))
	}

	return q
}

func (q *query) Find(res interface{}) (err error) {
	defer q.reset()

	resRt := reflect.TypeOf(res)
	if resRt.Kind() != reflect.Ptr {
		err = errs.ResIsNotPtr
		return
	}

	resRt = resRt.Elem()
	db := q.db.Model(reflect.New(resRt).Interface())
	if q.order != "" {
		db = db.Order(q.order)
	}
	if q.whereSql != "" {
		db = db.Where(q.whereSql, q.whereArgs...)
	}
	if q.page > 0 && q.pageSize > 0 {
		db = db.Offset((q.page - 1) * q.pageSize).Limit(q.pageSize)
	}
	if len(q.opts) > 0 {
		for _, o := range q.opts {
			if v, ok := o.(IOption); ok {
				db = v.Apply(db)
			}
		}
	}

	err = db.Find(res).Error

	return
}

func (q *query) First(res interface{}) (err error) {
	defer q.reset()
	db := q.db
	if q.order != "" {
		db = db.Order(q.order)
	}
	if q.whereSql != "" {
		db = db.Where(q.whereSql, q.whereArgs...)
	}
	if len(q.opts) > 0 {
		for _, o := range q.opts {
			if v, ok := o.(IOption); ok {
				db = v.Apply(db)
			}
		}
	}

	err = db.First(res).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (q *query) Asc(fields ...string) goresource.IQuery {
	if len(fields) > 0 {
		q.genOrder(" ASC", fields...)
	}

	return q
}

func (q *query) Desc(fields ...string) goresource.IQuery {
	if len(fields) > 0 {
		q.genOrder(" DESC", fields...)
	}

	return q
}

func (q *query) Page(page int) goresource.IQuery {
	if page > 0 {
		q.page = page
	}

	return q
}

func (q *query) PageSize(pageSize int) goresource.IQuery {
	if pageSize > 0 {
		q.pageSize = pageSize
	}

	return q
}

func (q *query) ToArray(res interface{}) (err error) {
	defer q.reset()

	return q.Find(res)
}

func (q *query) Where(args ...interface{}) goresource.IQuery {
	if len(args) > 0 {
		q.whereSql = args[0].(string)
		if len(args) > 1 {
			q.whereArgs = args[1:]
		}
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

func (q *query) genOrder(suffix string, fields ...string) {
	if q.order == "" {
		q.order = strings.Join(fields, ", ") + " " + suffix
	} else {
		q.order = q.order + ", " + strings.Join(fields, ", ") + " " + suffix
	}
}

func (q *query) reset() {
	q.page = 0
	q.pageSize = 0
	q.fields = make([]string, 0)
	q.order = ""
	q.whereArgs = make([]interface{}, 0)
}
