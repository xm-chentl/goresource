package goresource

import "github.com/xm-chentl/goresource/dbtype"

// IDbFactory 数据库实例工厂
type IFactory interface {
	BuildByType(dbtype.Value) (IResource, error) // 类型与数据库 一对一
	BuildByName(string) (IResource, error)       // 名称与数据库 一对一
}
