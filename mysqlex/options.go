package mysqlex

type SaveOptionBySelect struct {
	Fields []string
}

func (s SaveOptionBySelect) Desc() {}

type SaveOptionByOmit struct {
	Fields []string
}

func (s SaveOptionByOmit) Desc() {}
