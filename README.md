## 数据库工厂

面向业务逻辑代码，数据库操作统一接口协议

### 数据库实例工厂

```go
// IDbFactory 数据库实例工厂
type IDbFactory interface {
	BuildByType(dbtype.Value) (IFactory, error) // 类型与数据库 一对一
	BuildByName(string) (IFactory, error)       // 名称与数据 一对一
}
```

### 数据库工厂

```go
// IFactory 数据库实现
// Todo: Uow暂时未实现
type IFactory interface {
	Db(...interface{}) IRepository
	Uow() IUnitOfWork
}

```

### 使用示例

```go
```