package mongoex

import (
	"context"

	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/dbtype"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	ctx            context.Context
	database       *mongo.Database
	repositoryBase *goresource.RepositoryBase
	uow            *unitOfWork
}

func (r *repository) Create(entry goresource.IDbModel, args ...interface{}) (err error) {
	if v, ok := entry.GetID().(primitive.ObjectID); ok {
		if v.IsZero() {
			entry.SetID(primitive.NewObjectID())
		}
	}
	if r.uow != nil {
		r.uow.commitCreate(entry)
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.Mongo, r.uow)
		}

		return
	}

	result, err := r.database.Collection(entry.Table()).InsertOne(r.ctx, entry)
	if err != nil {
		return
	}
	if result.InsertedID != nil {
		entry.SetID(result.InsertedID)
	}

	return
}

// delete args 0 -> 支持many
func (r *repository) Delete(entry goresource.IDbModel, args ...interface{}) (err error) {
	if r.uow != nil {
		r.uow.commitDelete(entry, args...)
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.Mongo, r.uow)
		}

		return
	}
	if len(args) == 0 {
		_, err = r.database.Collection(entry.Table()).DeleteOne(r.ctx, bson.M{"_id": entry.GetID()})
	} else {
		_, err = r.database.Collection(entry.Table()).DeleteMany(r.ctx, args[0])
	}

	return
}

// Update args 0 upset 1 filter
func (r *repository) Update(entry goresource.IDbModel, args ...interface{}) (err error) {
	if r.uow != nil {
		r.uow.commitUpdate(entry, args...)
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.Mongo, r.uow)
		}

		return
	}

	collectionDb := r.database.Collection(entry.Table())
	filter := bson.M{"_id": entry.GetID()}
	// 默认更新完全
	if len(args) == 0 {
		_, err = collectionDb.UpdateOne(r.ctx, filter, bson.M{"$set": entry})
		return
	}
	// one
	if len(args) == 1 && args[0] != nil {
		_, err = collectionDb.UpdateOne(r.ctx, filter, args[0])
		return
	}
	// many
	if len(args) == 2 && args[0] != nil && args[1] != nil {
		_, err = collectionDb.UpdateMany(r.ctx, args[1], args[0])
		return
	}

	return
}

func (r *repository) Query() goresource.IQuery {
	return &query{
		ctx:      r.ctx,
		database: r.database,
		filter:   bson.M{},
		orders:   make([]string, 0),
		orderBy:  make([]string, 0),
		opts:     make([]IOption, 0),
	}
}
