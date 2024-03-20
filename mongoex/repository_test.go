package mongoex

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type testPerson struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func (t testPerson) Table() string {
	return "test-person"
}

func (t testPerson) GetID() interface{} {
	return t.ID
}

func (t *testPerson) SetID(v interface{}) {
	if idv, ok := v.(string); ok {
		t.ID = idv
	}
}

type modelIDString struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
}

func (m modelIDString) GetID() interface{} {
	return m.ID
}

func (m *modelIDString) SetID(v interface{}) {
	if idv, ok := v.(primitive.ObjectID); ok {
		m.ID = idv
	}
}

type objectEntry2 struct {
	modelIDString `bson:",inline"`

	Desc string `bson:"desc"`
}

func (m objectEntry2) Table() string {
	return "object_entry_2"
}

// type objectEntry3 struct {
// 	modelIDString `bson:",inline"`

// 	Desc string `bson:"desc"`
// }

// func (m objectEntry3) Table() string {
// 	return "object_entry_2"
// }

type objectEntry struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	Desc string `bson:"desc"`
}

func (m objectEntry) Table() string {
	return "object_entry"
}

func (m objectEntry) GetID() interface{} {
	return m.ID
}

func (m *objectEntry) SetID(v interface{}) {
	m.ID = v.(primitive.ObjectID)
}

func Test_repository_Create(test *testing.T) {
	repo, err := getClient()
	if err != nil {
		test.Fatal("err", err)
	}

	test.Run("success", func(t *testing.T) {
		addEntry := testPerson{
			ID:   "test-001",
			Name: "test-name-001",
		}
		if err = repo.Create(&addEntry); err != nil {
			t.Fatal("err", err)
		}

		db := repo.database.Collection(addEntry.Table())
		res := db.FindOne(context.Background(), bson.M{"_id": "test-001"})
		if err = res.Err(); err != nil {
			t.Fatal("err", err)
		}

		var entry testPerson
		if err = res.Decode(&entry); err != nil {
			t.Fatal("err", err)
		}

		a := assert.New(t)
		a.Equal(addEntry, entry)

		_, err = db.DeleteOne(context.Background(), bson.M{"_id": "test-001"})
		if err != nil {
			t.Fatal("err", err)
		}
	})

	test.Run("primitive.objectID.create", func(t *testing.T) {
		addEntry := objectEntry{
			Desc: "test_object_id",
		}
		if err = repo.Create(&addEntry); err != nil {
			t.Fatal(err, err)
		}

		db := repo.database.Collection(addEntry.Table())
		defer func() {
			_, _ = db.DeleteOne(context.Background(), bson.M{"_id": addEntry.ID})
		}()

		res := db.FindOne(context.Background(), bson.M{"_id": addEntry.ID})
		if res.Err() != nil {
			t.Fatal("err", res.Err())
		}

		queryEntry := objectEntry{}
		if err = res.Decode(&queryEntry); err != nil {
			t.Fatal("err", err)
		}

		a := assert.New(t)
		a.Equal(addEntry, queryEntry)
	})

	test.Run("primitive.objectID.create_inline", func(t *testing.T) {
		addEntry := objectEntry2{
			Desc: "test_object_id_2",
		}
		if err = repo.Create(&addEntry); err != nil {
			t.Fatal(err, err)
		}

		db := repo.database.Collection(addEntry.Table())
		defer func() {
			_, _ = db.DeleteOne(context.Background(), bson.M{"_id": addEntry.ID})
		}()

		res := db.FindOne(context.Background(), bson.M{"_id": addEntry.ID})
		if res.Err() != nil {
			t.Fatal("err", res.Err())
		}

		queryEntry := objectEntry2{}
		if err = res.Decode(&queryEntry); err != nil {
			t.Fatal("err", err)
		}
		if !addEntry.ID.IsZero() {
			t.Fatal("err")
		}
	})
}

func Test_repository_delete(test *testing.T) {
	repo, err := getClient()
	if err != nil {
		test.Fatal("err")
	}

	test.Run("deleteOne.success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "delete-test-001",
				Name: "delete-test",
			},
		}
		ctx := context.Background()
		db := repo.database.Collection(addEntries[0].Table())
		defer func() {
			for _, entry := range addEntries {
				_, err = db.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()

		var err error
		for index := range addEntries {
			entry := addEntries[index]
			_, err = db.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		deleteEntry := testPerson{
			ID: addEntries[0].ID,
		}
		if err = repo.Delete(&deleteEntry); err != nil {
			t.Fatal("err", err)
		}

		res := db.FindOne(ctx, bson.M{"_id": deleteEntry.ID})
		a := assert.New(t)
		a.Equal(mongo.ErrNoDocuments, res.Err())
	})

	test.Run("deleteMany", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "delete-test-002",
				Name: "delete-many-test",
			},
			{
				ID:   "delete-test-003",
				Name: "delete-many-test",
			},
		}
		ctx := context.Background()
		collectionDb := repo.database.Collection(addEntries[0].Table())
		defer func() {
			for _, entry := range addEntries {
				_, err = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()

		var err error
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}
		if err = repo.Delete(&addEntries[0], bson.M{"name": addEntries[0].Name}); err != nil {
			t.Fatal("err", err)
		}

		count, err := collectionDb.CountDocuments(ctx, bson.M{"name": addEntries[0].Name})
		if err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(int64(0), count)
	})

	test.Run("primitive.objectID.deleteOne.success", func(t *testing.T) {

	})
}

func Test_repository_Update(test *testing.T) {
	repo, err := getClient()
	if err != nil {
		test.Fatal("err")
	}

	test.Run("update.success", func(t *testing.T) {
		addEntry := testPerson{
			ID:   "update-test-001",
			Name: "update-test-001",
		}

		ctx := context.Background()
		db := repo.database.Collection(addEntry.Table())
		defer func() {
			_, err = db.DeleteOne(ctx, bson.M{"_id": addEntry.ID})
			if err != nil {
				t.Fatal("err")
			}
		}()

		_, err := db.InsertOne(ctx, &addEntry)
		if err != nil {
			t.Fatal("err")
		}

		addEntry.Name = "update-test"
		if err = repo.Update(&addEntry); err != nil {
			t.Fatal("err")
		}

		entry := testPerson{}
		res := db.FindOne(ctx, bson.M{"_id": addEntry.ID})
		if err = res.Decode(&entry); err != nil {
			t.Fatal("err")
		}
		a := assert.New(t)
		a.Equal(addEntry, entry)
	})

	test.Run("updateOne.upset.success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "update-test-002",
				Name: "update-test-002",
			},
			{
				ID:   "update-test-003",
				Name: "update-test-003",
			},
		}
		ctx := context.Background()
		db := repo.database.Collection(addEntries[0].Table())
		defer func() {
			for _, entry := range addEntries {
				_, err = db.DeleteOne(ctx, bson.M{"_id": entry.ID})
				if err != nil {
					t.Fatal("err")
				}
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err := db.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		updateEntry := testPerson{
			ID:   addEntries[1].ID,
			Name: "update-test-ctl",
		}
		err = repo.Update(&updateEntry, bson.M{"$set": bson.M{"name": updateEntry.Name}})
		if err != nil {
			t.Fatal("err")
		}

		entry := testPerson{}
		res := db.FindOne(ctx, bson.M{"_id": updateEntry.ID})
		if err = res.Decode(&entry); err != nil {
			t.Fatal("err")
		}
		a := assert.New(t)
		a.Equal(updateEntry, entry)
	})

	test.Run("updateMoney.upset.success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "update-many-test-002",
				Name: "update-test",
			},
			{
				ID:   "update-many-test-003",
				Name: "update-test",
			},
		}
		ctx := context.Background()
		db := repo.database.Collection(addEntries[0].Table())
		defer func() {
			for _, entry := range addEntries {
				_, err = db.DeleteOne(ctx, bson.M{"_id": entry.ID})
				if err != nil {
					t.Fatal("err")
				}
			}
		}()

		expectedEntries := make([]testPerson, 0)
		for index := range addEntries {
			entry := addEntries[index]
			_, err := db.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
			entry.Name = "update-test-ctl"
			expectedEntries = append(expectedEntries, entry)
		}

		err = repo.Update(&testPerson{}, bson.M{"$set": bson.M{"name": expectedEntries[0].Name}}, bson.M{"name": addEntries[0].Name})
		if err != nil {
			t.Fatal("err")
		}

		c, err := db.Find(ctx, bson.M{"name": expectedEntries[0].Name})
		if err != nil {
			t.Fatal("err")
		}
		defer c.Close(ctx)

		var entries []testPerson
		for c.Next(ctx) {
			entry := testPerson{}
			if err = c.Decode(&entry); err != nil {
				t.Fatal("err")
			}
			entries = append(entries, entry)
		}
		a := assert.New(t)
		a.Equal(2, len(entries))
		a.Equal(expectedEntries, entries)
	})
}
