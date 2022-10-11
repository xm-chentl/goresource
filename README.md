## 数据资源

面向业务逻辑代码，数据库资源统一接口协议

### 数据资源工厂

```go
// IFactory 数据库实例工厂
type IFactory interface {
	BuildByType(dbtype.Value) (IResource, error) // 类型与数据库 一对一
	BuildByName(string) (IResource, error)       // 名称与数据 一对一
}
```

### 数据库资源

```go
// IResource 数据资源实例
type IResource interface {
	Db(...interface{}) IRepository
	Uow() IUnitOfWork
}

```

### 使用示例

```go
```