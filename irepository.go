package goresource

type IRepository interface {
	Create(entry IDbModel, args ...interface{}) error
	Delete(entry IDbModel, args ...interface{}) error
	Update(entry IDbModel, args ...interface{}) error
	Query() IQuery
}
