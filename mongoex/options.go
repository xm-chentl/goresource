package mongoex

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type IOption interface {
	Apply(*mongo.Database) *mongo.Collection
}

type OptionCollection struct {
	Name string
}

func (o OptionCollection) Apply(db *mongo.Database) *mongo.Collection {
	return db.Collection(o.Name)
}

func NewOptionCollection(format string, args ...interface{}) IOption {
	return &OptionCollection{
		Name: fmt.Sprintf(format, args...),
	}
}
