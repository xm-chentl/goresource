package goresource

// IFactory 数据库实现
type IResource interface {
	Db(...interface{}) IRepository
	Uow() IUnitOfWork
}
