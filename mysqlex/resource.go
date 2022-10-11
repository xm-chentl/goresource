package mysqlex

import (
	"context"

	"github.com/xm-chentl/goresource"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type resource struct {
	dsn string
	db  *gorm.DB
}

// Db 参数请按 ctx uow
func (f *resource) Db(args ...interface{}) goresource.IRepository {
	repo := &repository{}
	for _, a := range args {
		if ctx, ok := a.(context.Context); ok {
			repo.db = f.db.WithContext(ctx)
		} else if uow, ok := a.(*unitOfWork); ok {
			repo.uow = uow
		} else if uow, ok := a.(goresource.IUnitOfWork); ok {
			repo.uow = newUnitOfWork(f.db)
			repo.repositoryBase = goresource.NewRepository(uow)
		}
	}
	if repo.db == nil {
		repo.db = f.db
	}

	return repo
}

func (f *resource) Uow() goresource.IUnitOfWork {
	return newUnitOfWork(f.db)
}

func New(dsn string) goresource.IResource {
	if dsn == "" {
		panic("mysqlex.New parameter dsn is empty")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("mysqlex.New open database is failed err: " + err.Error())
	}
	return &resource{
		dsn: dsn,
		db:  db,
	}
}
