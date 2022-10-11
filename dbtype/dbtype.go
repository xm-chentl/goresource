package dbtype

// Value 数据库类型
type Value string

func (v Value) String() string {
	return string(v)
}

const (
	Mongo     Value = "mongo"
	MySQL     Value = "mysql"
	TimeScale Value = "timescale"
)
