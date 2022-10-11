package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type pool struct {
	ctx     context.Context
	pgxPool *pgxpool.Pool
}

func (p *pool) getConn() (conn *pgxpool.Conn, err error) {
	conn, err = p.pgxPool.Acquire(p.ctx)
	if err != nil {
		return
	}

	return
}
