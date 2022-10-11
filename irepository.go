package goresource

type IRepository interface {
	Create(entry IDbModel) error
	Delete(entry IDbModel, args ...interface{}) error
	Update(entry IDbModel, args ...interface{}) error
	Query() IQuery
}
