package goresource

type IQuery interface {
	Count(entry IDbModel) (count int64, err error)
	Exec(res interface{}, args ...interface{}) (err error)
	// Fields 指定筛选字段（根据各资源的使用方式 如: mysql "field1", "field2" mongo bson.M{"field1":1, "field2":1}）
	Fields(args ...interface{}) IQuery
	Find(res interface{}, opts ...IFindOptions) error
	First(res interface{}) error
	Asc(fields ...string) IQuery
	Desc(fields ...string) IQuery
	Page(page int) IQuery
	PageSize(pageSize int) IQuery
	ToArray(res interface{}) error
	Where(args ...interface{}) IQuery
}
