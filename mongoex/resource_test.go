package mongoex

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testConnStr = "mongodb://localhost:27017"
)

func Test_Factory_NewFactory(test *testing.T) {
	test.Run("Connect.Success", func(t *testing.T) {
		_, err := getClient()
		if err != nil {
			t.Fatal("err")
		}
	})
}

func getClient() (repo repository, err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(testConnStr))
	if err != nil {
		return
	}
	if err = client.Connect(context.Background()); err != nil {
		return
	}
	repo = repository{
		database: client.Database("testdb"),
	}

	return
}
