package mysqlex

import (
	"fmt"

	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/dbtype"
	"github.com/xm-chentl/goresource/repositorytype"

	"gorm.io/gorm"
)

type repository struct {
	db             *gorm.DB
	uow            *unitOfWork
	repositoryBase *goresource.RepositoryBase
}

func (r repository) Create(entry goresource.IDbModel) (err error) {
	if r.uow != nil {
		r.uow.commitQueues = append(r.uow.commitQueues, commitQueueItem{
			rt:    repositorytype.Create,
			entry: entry,
		})
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.MySQL, r.uow)
		}
		return
	}

	err = r.db.Model(entry).Create(entry).Error
	fmt.Println(">>>", entry.GetID())

	return
}

func (r repository) Delete(entry goresource.IDbModel, args ...interface{}) (err error) {
	if r.uow != nil {
		r.uow.commitQueues = append(r.uow.commitQueues, commitQueueItem{
			rt:    repositorytype.Delete,
			entry: entry,
			args:  args,
		})
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.MySQL, r.uow)
		}
		return
	}

	err = r.db.Model(entry).Delete(entry, args...).Error
	return
}

func (r repository) Update(entry goresource.IDbModel, args ...interface{}) (err error) {
	if r.uow != nil {
		r.uow.commitQueues = append(r.uow.commitQueues, commitQueueItem{
			rt:    repositorytype.Update,
			entry: entry,
			args:  args,
		})
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.MySQL, r.uow)
		}

		return
	}

	db := r.db.Model(entry)
	if len(args) > 0 {
		for _, a := range args {
			if v, ok := a.(SaveOptionByOmit); ok {
				db.Omit(v.Fields...)
			}
		}
	}

	return db.Updates(entry).Error
}

func (r repository) Query() goresource.IQuery {
	return &query{
		db: r.db,
	}
}
