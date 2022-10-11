package metadata

type ITable interface {
	AutoIncrementColumn() IColumn
	Name() string
	Columns() []IColumn
	ColumnMap() map[string]IColumn
	PrimaryKeyColumn() IColumn
}
