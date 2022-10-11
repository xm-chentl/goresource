package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xm-chentl/goresource"
)

type resource struct {
	dsn     string
	pgxPool *pgxpool.Pool
}

// Todo: 联合事务有问题
func (f resource) Db(args ...interface{}) goresource.IRepository {
	repo := &repository{
		pool: &pool{
			pgxPool: f.pgxPool,
		},
	}
	for index := range args {
		if ctx, ok := args[index].(context.Context); ok {
			repo.ctx = ctx
			continue
		}
		if uow, ok := args[index].(*unitOfWork); ok {
			repo.uow = uow
			continue
		}
		if uow, ok := args[index].(goresource.IUnitOfWork); ok {
			repo.uow = &unitOfWork{
				pool:          repo.pool,
				addOfQueue:    make([]commitQueueInfo, 0),
				deleteOfQueue: make([]commitQueueInfo, 0),
				updateOfQueue: make([]commitQueueInfo, 0),
			}
			repo.repositoryBase = goresource.NewRepository(uow)
		}
	}
	if repo.ctx == nil {
		repo.ctx = context.Background()
	}
	repo.pool.ctx = repo.ctx
	if repo.uow != nil {
		repo.uow.ctx = repo.ctx
	}

	return repo
}

func (f resource) Uow() goresource.IUnitOfWork {
	return &unitOfWork{
		pool: &pool{
			pgxPool: f.pgxPool,
		},
		addOfQueue:    make([]commitQueueInfo, 0),
		updateOfQueue: make([]commitQueueInfo, 0),
	}
}

func New(dsn string) goresource.IResource {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic("connect to database config faild: " + err.Error())
	}

	ctx := context.Background()
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		panic("Unable to connect to database: " + dsn)
	}
	if err = pool.Ping(ctx); err != nil {
		panic("connect to database faild: " + err.Error())
	}

	return &resource{
		pgxPool: pool,
		dsn:     dsn,
	}
}

func NewByGorm(connStr string) goresource.IFactory {

	return nil
}
