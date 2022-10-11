package repositorytype

type Value int

const (
	Create Value = iota
	Delete
	Update
)
