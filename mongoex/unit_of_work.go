package mongoex

import (
	"context"
	"sync"
	"time"

	"github.com/xm-chentl/goresource"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type commitQueueInfo struct {
	entry goresource.IDbModel
	args  []interface{}
}

type unitOfWork struct {
	ctx      context.Context
	database *mongo.Database

	collectionMap sync.Map
	createQueue   []commitQueueInfo
	deleteQueue   []commitQueueInfo
	updateQueue   []commitQueueInfo
}

func (u *unitOfWork) Commit() (err error) {
	defer u.reset()

	var collectionDb *mongo.Collection
	for index := range u.createQueue {
		item := u.createQueue[index]
		if v, ok := item.entry.GetID().(primitive.ObjectID); ok {
			if v.IsZero() {
				item.entry.SetID(primitive.NewObjectID())
			}
		}

		collectionDb = u.getCollection(item.entry)
		_, err = collectionDb.InsertOne(u.ctx, item.entry)
		if err != nil {
			return
		}
	}
	for index := range u.deleteQueue {
		item := u.deleteQueue[index]
		collectionDb = u.getCollection(item.entry)
		if item.args != nil {
			_, err = collectionDb.DeleteOne(u.ctx, item.args)
		} else {
			_, err = collectionDb.DeleteOne(u.ctx, bson.M{"_id": item.entry.GetID()})
		}
		if err != nil {
			return
		}
	}
	for index := range u.updateQueue {
		item := u.updateQueue[index]
		collectionDb = u.getCollection(item.entry)
		filter := bson.M{"_id": item.entry.GetID()}
		// 默认更新完全
		if len(item.args) == 0 {
			_, err = collectionDb.UpdateOne(u.ctx, filter, bson.M{"$set": item.entry})
			continue
		}
		// one
		if len(item.args) == 1 && item.args[0] != nil {
			_, err = collectionDb.UpdateOne(u.ctx, filter, item.args[0])
			continue
		}
		// many
		if len(item.args) == 2 && item.args[0] != nil && item.args[1] != nil {
			_, err = collectionDb.UpdateMany(u.ctx, item.args[1], item.args[0])
			continue
		}
		if err != nil {
			return
		}
	}

	return
}

func (u *unitOfWork) commitCreate(entry goresource.IDbModel) {
	u.createQueue = append(u.createQueue, commitQueueInfo{
		entry: entry,
	})
}

func (u *unitOfWork) commitDelete(entry goresource.IDbModel, args ...interface{}) {
	u.deleteQueue = append(u.deleteQueue, commitQueueInfo{
		entry: entry,
		args:  args,
	})
}

func (u *unitOfWork) commitUpdate(entry goresource.IDbModel, args ...interface{}) {
	u.updateQueue = append(u.updateQueue, commitQueueInfo{
		entry: entry,
		args:  args,
	})
}

func (u *unitOfWork) getCollection(entry goresource.IDbModel) (collectionDb *mongo.Collection) {
	value, ok := u.collectionMap.Load(entry.Table())
	if !ok {
		value = u.database.Collection(entry.Table())
		u.collectionMap.Store(entry.Table(), value)
	}

	collectionDb = value.(*mongo.Collection)
	return
}

func (u *unitOfWork) reset() {
	u.createQueue = make([]commitQueueInfo, 0)
	u.deleteQueue = make([]commitQueueInfo, 0)
	u.updateQueue = make([]commitQueueInfo, 0)
	u.collectionMap.Range(func(key, _ interface{}) bool {
		u.collectionMap.Delete(key)
		return true
	})
}

func (u *unitOfWork) commit1() error {
	// 暂时不使用事务
	return u.database.Client().UseSession(u.ctx, func(sessionCtx mongo.SessionContext) (err error) {
		if err = sessionCtx.StartTransaction(); err != nil {
			return
		}

		var collectionDb *mongo.Collection
		for index := range u.createQueue {
			item := u.createQueue[index]
			collectionDb = u.getCollection(item.entry)
			_, err = collectionDb.InsertOne(sessionCtx, item.entry)
			if err != nil {
				return
			}
		}
		for index := range u.deleteQueue {
			item := u.deleteQueue[index]
			collectionDb = u.getCollection(item.entry)
			if item.args != nil {
				_, err = collectionDb.DeleteOne(sessionCtx, item.args)
			} else {
				_, err = collectionDb.DeleteOne(sessionCtx, bson.M{"_id": item.entry.GetID()})
			}
			if err != nil {
				return
			}
		}
		for index := range u.updateQueue {
			item := u.updateQueue[index]
			collectionDb = u.getCollection(item.entry)
			if item.args != nil {
				_, err = collectionDb.UpdateOne(sessionCtx, bson.M{"_id": item.entry.GetID()}, item.args)
			} else {
				_, err = collectionDb.UpdateOne(sessionCtx, bson.M{"_id": item.entry.GetID()}, bson.M{"$set": item.entry})
			}
			if err != nil {
				return
			}
		}
		err = sessionCtx.CommitTransaction(context.Background())

		return
	})
}

// todo: 副本集使用方式
func (u *unitOfWork) commit2() (err error) {
	defer u.reset()

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	session, err := u.database.Client().StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(context.TODO())

	sessionCtx := mongo.NewSessionContext(ctx, session)
	if err = session.StartTransaction(); err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = session.AbortTransaction(context.Background())
		}
	}()

	var collectionDb *mongo.Collection
	for index := range u.createQueue {
		item := u.createQueue[index]
		collectionDb = u.getCollection(item.entry)
		_, err = collectionDb.InsertOne(sessionCtx, item.entry)
		if err != nil {
			return
		}
	}
	for index := range u.deleteQueue {
		item := u.deleteQueue[index]
		collectionDb = u.getCollection(item.entry)
		if item.args != nil {
			_, err = collectionDb.DeleteOne(sessionCtx, item.args)
		} else {
			_, err = collectionDb.DeleteOne(sessionCtx, bson.M{"_id": item.entry.GetID()})
		}
		if err != nil {
			return
		}
	}
	for index := range u.updateQueue {
		item := u.updateQueue[index]
		collectionDb = u.getCollection(item.entry)
		if item.args != nil {
			_, err = collectionDb.UpdateOne(sessionCtx, bson.M{"_id": item.entry.GetID()}, item.args)
		} else {
			_, err = collectionDb.UpdateOne(sessionCtx, bson.M{"_id": item.entry.GetID()}, bson.M{"$set": item.entry})
		}
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}

	return
}

func newUnitOfWork(database *mongo.Database) *unitOfWork {
	return &unitOfWork{
		database:    database,
		createQueue: make([]commitQueueInfo, 0),
		deleteQueue: make([]commitQueueInfo, 0),
		updateQueue: make([]commitQueueInfo, 0),
	}
}
