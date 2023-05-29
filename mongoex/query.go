package mongoex

import (
	"context"
	"reflect"

	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/errs"
	"github.com/xm-chentl/goresource/tools"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type query struct {
	ctx        context.Context
	database   *mongo.Database
	projection interface{}
	filter     bson.M
	page       int
	pageSize   int
	orders     []string // 1
	orderBy    []string // -1
	opts       []IOption
}

func (q *query) Asc(fields ...string) goresource.IQuery {
	if len(fields) > 0 {
		filterFields := make([]string, 0)
		for index := range fields {
			field := fields[index]
			if field != "" {
				filterFields = append(filterFields, field)
			}
		}
		q.orders = append(q.orders, filterFields...)
	}
	return q
}

func (q *query) Count(entry goresource.IDbModel) (res int64, err error) {
	res, err = q.database.Collection(entry.Table()).CountDocuments(q.ctx, q.filter)

	return
}

func (q *query) Desc(fields ...string) goresource.IQuery {
	if len(fields) > 0 {
		filterFields := make([]string, 0)
		for index := range fields {
			field := fields[index]
			if field != "" {
				filterFields = append(filterFields, field)
			}
		}
		q.orderBy = append(q.orderBy, filterFields...)
	}
	return q
}

func (q query) Exec(res interface{}, args ...interface{}) (err error) {
	return
}

func (q *query) Fields(args ...interface{}) goresource.IQuery {
	if len(args) > 0 {
		q.projection = args[0]
	}

	return q
}

func (q *query) Find(res interface{}) (err error) {
	defer q.reset()

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

	var collectionDb *mongo.Collection
	for _, opt := range q.opts {
		collectionDb = opt.Apply(q.database)
	}
	if collectionDb == nil {
		newEntry := reflect.New(resRt).Interface().(goresource.IDbModel)
		collectionDb = q.database.Collection(newEntry.Table())
	}

	opt := &options.FindOptions{}
	if q.page > 0 || q.pageSize > 0 {
		opt.SetSkip(int64((q.page - 1) * q.pageSize)).SetLimit(int64(q.pageSize))
	}
	if len(q.orders) > 0 || len(q.orderBy) > 0 {
		sort := make(bson.D, 0)
		for index := range q.orders {
			sort = append(sort, bson.E{
				Key:   q.orders[index],
				Value: 1,
			})
		}
		for index := range q.orderBy {
			sort = append(sort, bson.E{
				Key:   q.orderBy[index],
				Value: -1,
			})
		}
		opt.SetSort(sort)
	}
	if q.projection != nil {
		opt.SetProjection(q.projection)
	}

	cursor, err := collectionDb.Find(q.ctx, q.filter, opt)
	if err != nil {
		return
	}

	tempSlice := reflect.MakeSlice(reflect.TypeOf(res).Elem(), 0, 0)
	for cursor.Next(q.ctx) {
		mappingInst := reflect.New(resRt).Interface()
		err := cursor.Decode(mappingInst)
		if err != nil {
			break
		}
		tempSlice = reflect.Append(tempSlice, reflect.ValueOf(mappingInst).Elem())
	}
	resRv.Elem().Set(tempSlice)

	return
}

func (q *query) First(res interface{}) (err error) {
	defer q.reset()

	entry, ok := res.(goresource.IDbModel)
	if !ok {
		err = errs.ResIsNotIDbModel
		return
	}
	if len(q.filter) == 0 && !tools.IsEmpty(entry.GetID()) {
		q.filter = bson.M{"_id": entry.GetID()}
	}

	opt := &options.FindOneOptions{}
	if len(q.orders) > 0 || len(q.orderBy) > 0 {
		sort := make(bson.D, 0)
		for index := range q.orders {
			sort = append(sort, bson.E{
				Key:   q.orders[index],
				Value: 1,
			})
		}
		for index := range q.orderBy {
			sort = append(sort, bson.E{
				Key:   q.orderBy[index],
				Value: -1,
			})
		}
		opt.SetSort(sort)
	}
	if q.projection != nil {
		opt.SetProjection(q.projection)
	}

	collectionDb := q.database.Collection(entry.Table())
	result := collectionDb.FindOne(q.ctx, q.filter, opt)
	err = result.Err()
	if err == mongo.ErrNoDocuments {
		err = nil
		return
	}
	if err != nil {
		return
	}
	err = result.Decode(entry)

	return
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
	if q.pageSize == 0 {
		q.pageSize = 10
	}

	return q
}

func (q query) ToArray(res interface{}) error {
	return q.Find(res)
}

func (q *query) Where(args ...interface{}) goresource.IQuery {
	if len(args) > 0 {
		q.filter = args[0].(bson.M)
	}

	return q
}

func (q *query) SetOpts(opts ...interface{}) goresource.IQuery {
	return q
}

func (q *query) reset() {
	q.filter = bson.M{}
	q.projection = nil
}
