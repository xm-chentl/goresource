package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	connStr = "postgres://test:123456@47.98.248.82:15432/test?pool_min_conns=10&pool_max_conns=100"
)

func Test_NewFactory(test *testing.T) {
	test.Run("Connect Success", func(t *testing.T) {
		repo, err := getRepo()
		if err != nil {
			t.Fatal(err)
		}
		defer repo.pool.pgxPool.Close()

		_, err = repo.pool.pgxPool.Exec(repo.ctx, createTestPesonSql)
		if err != nil {
			t.Fatal("init create test_person faild: ", err)
		}

		_, err = repo.pool.pgxPool.Exec(repo.ctx, createTimeStructSql)
		if err != nil {
			t.Fatal("init create test_time faild: ", err)
		}
	})
}

func getRepo() (repo repository, err error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic("connect to database config faild: " + err.Error())
	}

	pgxPool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		panic("Unable to connect to database: " + connStr)
	}
	repo = repository{
		pool: &pool{
			pgxPool: pgxPool,
			ctx:     ctx,
		},
		ctx: ctx,
	}

	return
}
