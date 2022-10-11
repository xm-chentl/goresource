package postgres

import (
	"context"
	"fmt"

	"github.com/xm-chentl/goresource"
	"github.com/xm-chentl/goresource/dbtype"
	"github.com/xm-chentl/goresource/errs"
	"github.com/xm-chentl/goresource/postgres/grammar"
	"github.com/xm-chentl/goresource/postgres/metadata"
	"github.com/xm-chentl/goresource/tools"
)

type repository struct {
	ctx            context.Context
	repositoryBase *goresource.RepositoryBase
	pool           *pool
	uow            *unitOfWork
}

func (r *repository) Create(entry goresource.IDbModel) (err error) {
	sql, args := grammar.Insert(metadata.Get(entry), entry)
	if r.uow != nil {
		r.uow.addQueue(sql, args...)
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.TimeScale, r.uow)
		}
		return
	}

	err = r.exec(sql, args...)

	return
}

// args 0 filter
func (r repository) Delete(entry goresource.IDbModel, args ...interface{}) (err error) {
	newArgs := make([]interface{}, 0)
	newArgs = append(newArgs, args...)
	table := metadata.Get(entry)
	pkColumn := table.PrimaryKeyColumn()
	// 没筛选条件 && 存在主键 && 主键有值 默认是id
	if len(newArgs) == 0 && pkColumn != nil && !tools.IsEmpty(entry.GetID()) {
		newArgs = append(newArgs, fmt.Sprintf("%s = $1", pkColumn.Field()))
		newArgs = append(newArgs, entry.GetID())
	}
	if len(newArgs) == 0 {
		// 不允许全量删除
		err = errs.DeleteFullNotAllowed
		return
	}

	sql, args := grammar.Delete(metadata.Get(entry), entry, args...)
	if r.uow != nil {
		r.uow.deleteQueue(sql, args...)
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.TimeScale, r.uow)
		}
		return
	}
	err = r.exec(sql, args...)

	return
}

// args 0 update-fields 1 filter (0 where-sql 1 where-args)
func (r repository) Update(entry goresource.IDbModel, args ...interface{}) (err error) {
	var updateFields []string
	var ok bool
	if len(args) > 0 {
		updateFields, ok = args[0].([]string)
		if !ok {
			err = ErrUpdateSetArgsIsNotArray
			return
		}
	}

	newArgs := make([]interface{}, 0)
	var where string
	if len(args) > 1 {
		where, ok = args[1].(string)
		if !ok {
			err = ErrUpdateSetQueryIsNotString
			return
		}
		newArgs = append(newArgs, where)
	}
	if len(args) > 2 {
		newArgs = append(newArgs, args[2:]...)
	}

	table := metadata.Get(entry)
	pkColumn := table.PrimaryKeyColumn()
	if len(newArgs) == 0 && pkColumn != nil {
		argsCount := len(table.Columns())
		if len(updateFields) > 0 {
			argsCount = len(updateFields)
		}
		newArgs = append(newArgs, fmt.Sprintf("%s = $%d", pkColumn.Field(), argsCount+1))
		newArgs = append(newArgs, entry.GetID())
	}
	if len(newArgs) == 0 {
		// 不允许全量更新（查询没有条件、主键时）
		err = errs.UpdateFullNotAllowed
		return
	}

	sql, args := grammar.Update(metadata.Get(entry), entry, updateFields, newArgs...)
	if r.uow != nil {
		r.uow.updateQueue(sql, args...)
		if r.repositoryBase != nil {
			r.repositoryBase.SetUow(dbtype.TimeScale, r.uow)
		}
		return
	}
	err = r.exec(sql, args...)

	return
}

func (r repository) exec(sql string, args ...interface{}) (err error) {
	conn, err := r.pool.getConn()
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Conn().Exec(r.ctx, sql, args...)

	return
}

func (r *repository) Query() goresource.IQuery {
	return &query{
		ctx:       r.ctx,
		pool:      r.pool,
		fields:    make([]string, 0),
		whereArgs: make([]interface{}, 0),
		orders:    make([]string, 0),
		orderBys:  make([]string, 0),
	}
}
