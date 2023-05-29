package mysqlex

import (
	"reflect"

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

func (r repository) Create(entry goresource.IDbModel, args ...interface{}) (err error) {
	whereArgs, opts, db, hook := optionApply(r.db, entry, args...)
	if r.uow != nil {
		r.uow.commitQueues = append(r.uow.commitQueues, commitQueueItem{
			rt:     repositorytype.Create,
			entry:  entry,
			args:   whereArgs,
			opts:   opts,
			filter: hook,
		})
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.MySQL, r.uow)
		}
		return
	}
	err = db.Model(entry).Create(entry).Error

	return
}

func (r repository) Delete(entry goresource.IDbModel, args ...interface{}) (err error) {
	whereArgs, opts, db, hook := optionApply(r.db, entry, args...)
	if r.uow != nil {
		r.uow.commitQueues = append(r.uow.commitQueues, commitQueueItem{
			rt:     repositorytype.Delete,
			entry:  entry,
			args:   whereArgs,
			opts:   opts,
			filter: hook,
		})
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.MySQL, r.uow)
		}
		return
	}
	err = db.Model(entry).Delete(entry, args...).Error

	return
}

func (r repository) Update(entry goresource.IDbModel, args ...interface{}) (err error) {
	whereArgs, opts, db, hook := optionApply(r.db, entry, args...)
	if r.uow != nil {
		r.uow.commitQueues = append(r.uow.commitQueues, commitQueueItem{
			rt:     repositorytype.Update,
			entry:  entry,
			args:   whereArgs,
			opts:   opts,
			filter: hook,
		})
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.MySQL, r.uow)
		}

		return
	}
	err = db.Updates(entry).Error

	return
}

func (r repository) Query() goresource.IQuery {
	return &query{
		db: r.db,
	}
}

func optionApply(db *gorm.DB, entry goresource.IDbModel, vs ...interface{}) (
	args []interface{},
	opts []IOption,
	dbRes *gorm.DB,
	hook HookFilter,
) {
	args = make([]interface{}, 0)
	opts = make([]IOption, 0)
	for _, v := range vs {
		if o, ok := v.(IOption); ok {
			opts = append(opts, o)
			rt := reflect.TypeOf(o)
			if rt.Kind() == reflect.Ptr {
				rt = rt.Elem()
			}
			switch rt {
			case reflect.TypeOf(OptionTableSuffix{}):
				// 创建表(存在则不重复创建)
				// todo: 每次都会判断，影响io、加上缓存(服务半掉再开起来，还是得判断至少一次)也会影响io，解决方案: 创建运营商和授权游戏的时候创建
				// todo: 更建议由业务端去做处理，把建表的逻辑开放出去
				mg := db.Migrator()
				ot := o.(*OptionTableSuffix)
				if !mg.HasTable(ot.Value) {
					_ = mg.AutoMigrate(entry)
					_ = mg.RenameTable(entry.Table(), ot.Value)
				}
			}
			db = o.Apply(db)
			continue
		}
		if h, ok := v.(HookFilter); ok {
			hook = h
			continue
		}
		args = append(args, v)
	}
	dbRes = db

	return
}
