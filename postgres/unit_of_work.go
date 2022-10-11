package postgres

import (
	"context"
)

type commitQueueInfo struct {
	sql  string
	args []interface{}
}

type unitOfWork struct {
	ctx  context.Context
	pool *pool

	addOfQueue    []commitQueueInfo
	updateOfQueue []commitQueueInfo
	deleteOfQueue []commitQueueInfo
}

func (u *unitOfWork) Commit() (err error) {
	defer u.reset()

	conn, err := u.pool.getConn()
	if err != nil {
		return
	}
	defer conn.Release()

	for index := range u.addOfQueue {
		item := u.addOfQueue[index]
		if _, err = conn.Exec(u.ctx, item.sql, item.args...); err != nil {
			return
		}
	}
	for index := range u.updateOfQueue {
		item := u.updateOfQueue[index]
		if _, err = conn.Exec(u.ctx, item.sql, item.args...); err != nil {
			return
		}
	}
	for index := range u.deleteOfQueue {
		item := u.deleteOfQueue[index]
		if _, err = conn.Exec(u.ctx, item.sql, item.args...); err != nil {
			return
		}
	}

	return
}

func (u *unitOfWork) Commit1() (err error) {
	conn, err := u.pool.getConn()
	if err != nil {
		return
	}
	defer conn.Release()

	tx, err := conn.Begin(u.ctx)
	if err != nil {
		return
	}
	defer func() {
		if tx != nil {
			if err != nil {
				_ = tx.Rollback(u.ctx)
			} else {
				_ = tx.Commit(u.ctx)
			}
		}
	}()

	for index := range u.addOfQueue {
		item := u.addOfQueue[index]
		if _, err = tx.Exec(u.ctx, item.sql, item.args...); err != nil {
			return
		}
	}

	return
}

func (u *unitOfWork) addQueue(sql string, args ...interface{}) {
	u.addOfQueue = append(u.addOfQueue, commitQueueInfo{
		sql:  sql,
		args: args,
	})
}

func (u *unitOfWork) updateQueue(sql string, args ...interface{}) {
	u.updateOfQueue = append(u.updateOfQueue, commitQueueInfo{
		sql:  sql,
		args: args,
	})
}

func (u *unitOfWork) deleteQueue(sql string, args ...interface{}) {
	u.deleteOfQueue = append(u.deleteOfQueue, commitQueueInfo{
		sql:  sql,
		args: args,
	})
}

func (u *unitOfWork) reset() {
	u.addOfQueue = make([]commitQueueInfo, 0)
	u.updateOfQueue = make([]commitQueueInfo, 0)
	u.deleteOfQueue = make([]commitQueueInfo, 0)
}
