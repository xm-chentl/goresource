package grammar

const (
	VarTag = "@"
)

type IInsert interface{}

type IDelete interface{}

type IUpdate interface{}

type ISelect interface{}

type IWhere interface {
	Asc(fields ...string)
	Desc(fields ...string)
	Page(page int)
	PageSize(size int)
	Where(args ...interface{})
}
