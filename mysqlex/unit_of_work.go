package mysqlex

import (
	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/repositorytype"

	"gorm.io/gorm"
)

type HookFilter func(entry goresource.IDbModel) goresource.IDbModel

func NewHookFilter(h func(goresource.IDbModel) goresource.IDbModel) HookFilter {
	return h
}

type commitQueueItem struct {
	rt     repositorytype.Value
	args   []interface{}
	opts   []IOption
	filter HookFilter
	entry  goresource.IDbModel
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
			isPointTable := false
			for _, o := range item.opts {
				if _, ok := o.(*OptionTableSuffix); ok {
					isPointTable = true
				}
				tx = o.Apply(tx)
			}
			// todo: 坑 (指定了.Table()，会覆盖整笨tx对象，没有清空)【暂时处理】
			if len(item.opts) == 0 || !isPointTable {
				tx.Table(item.entry.Table())
			}
			if item.filter != nil {
				item.entry = item.filter(item.entry)
			}
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
				if txErr = tx.Model(item.entry).Save(item.entry).Error; txErr != nil {
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
