package mysqlex

import (
	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/repositorytype"

	"gorm.io/gorm"
)

type commitQueueItem struct {
	rt    repositorytype.Value
	args  []interface{}
	opts  []interface{}
	entry goresource.IDbModel
}

type unitOfWork struct {
	db           *gorm.DB
	commitQueues []commitQueueItem
}

func (u unitOfWork) Commit() (err error) {
	if len(u.commitQueues) == 0 {
		return
	}

	err = u.db.Transaction(func(tx *gorm.DB) (txErr error) {
		for _, item := range u.commitQueues {
			if item.rt == repositorytype.Create {
				if txErr = tx.Model(item.entry).Create(item.entry).Error; txErr != nil {
					return
				}
			} else if item.rt == repositorytype.Delete {
				args := item.args
				if args == nil {
					args = make([]interface{}, 0)
				}
				if txErr = tx.Model(item.entry).Delete(item.entry, args...).Error; txErr != nil {
					return
				}
			} else if item.rt == repositorytype.Update {
				if len(item.args) > 0 {
					for _, a := range item.args {
						if v, ok := a.(SaveOptionByOmit); ok {
							tx = tx.Omit(v.Fields...)
						}
					}
				}
				if txErr = tx.Model(item.entry).Updates(item.entry).Error; txErr != nil {
					return
				}
			}
		}

		return
	})

	return
}

func newUnitOfWork(db *gorm.DB) *unitOfWork {
	return &unitOfWork{
		db:           db,
		commitQueues: make([]commitQueueItem, 0),
	}
}
