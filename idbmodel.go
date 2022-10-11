package goresource

// IDbModel 数据库模型
type IDbModel interface {
	GetID() interface{}
	Table() string
	SetID(v interface{})
}
