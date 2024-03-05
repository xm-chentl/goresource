package mongoex

import (
	"context"

	"github.com/xm-chentl/goresource"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type resource struct {
	dbName   string
	database *mongo.Database
}

func (f resource) Db(args ...interface{}) goresource.IRepository {
	repo := &repository{
		database: f.database,
	}
	for index := range args {
		if ctx, ok := args[index].(context.Context); ok {
			repo.ctx = ctx
		} else if uow, ok := args[index].(*unitOfWork); ok {
			repo.uow = uow
		} else if uow, ok := args[index].(goresource.IUnitOfWork); ok {
			repo.uow = newUnitOfWork(f.database)
			repo.repositoryBase = goresource.NewRepository(uow)
		}
	}
	if repo.ctx == nil {
		repo.ctx = context.Background()
	}
	if repo.uow != nil {
		repo.uow.ctx = repo.ctx
	}

	return repo
}

func (f resource) Uow() goresource.IUnitOfWork {
	return newUnitOfWork(f.database)
}

func New(dbName, dsn string) goresource.IResource {
	opt := options.Client().ApplyURI(dsn)
	client, err := mongo.NewClient(opt)
	if err != nil {
		panic("create connect to mongo faild err: " + err.Error())
	}
	if err = client.Connect(context.Background()); err != nil {
		panic("connect to mongo faild err: " + err.Error())
	}

	return &resource{
		dbName:   dbName,
		database: client.Database(dbName),
	}
}
