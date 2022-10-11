package goresource

type IUnitOfWork interface {
	Commit() error
}
