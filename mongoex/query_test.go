package mongoex

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_query_Asc(test *testing.T) {
	test.Run("filter.field", func(t *testing.T) {
		q := query{}
		q.Asc("")
		if len(q.orders) != 0 {
			t.Fatal("err")
		}
	})

	test.Run("success", func(t *testing.T) {
		q := query{}
		q.Asc("name")
		if len(q.orders) != 1 {
			t.Fatal("err")
		}
	})
}

func Test_query_Count(test *testing.T) {
	repo, err := getClient()
	if err != nil {
		test.Fatal("err", err)
	}
	test.Run("all", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "query-test-001",
				Name: "query-test-name-001",
			},
			{
				ID:   "query-test-002",
				Name: "query-test-name-002",
			},
			{
				ID:   "query-test-003",
				Name: "query-test-name-003",
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		count, err := repo.Query().Count(&testPerson{})
		if err != nil || count != 3 {
			t.Fatal("err", err)
		}
	})

	test.Run("filter", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "query-test-001",
				Name: "query-test-name-001",
			},
			{
				ID:   "query-test-002",
				Name: "query-test-name",
			},
			{
				ID:   "query-test-003",
				Name: "query-test-name",
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		count, err := repo.Query().Where(bson.M{"name": addEntries[1].Name}).Count(&testPerson{})
		if err != nil || count != 2 {
			t.Fatal("err", err)
		}
	})
}

func Test_query_Desc(test *testing.T) {
	test.Run("filter.field", func(t *testing.T) {
		q := query{}
		q.Desc("")
		if len(q.orderBy) != 0 {
			t.Fatal("err")
		}
	})

	test.Run("success", func(t *testing.T) {
		q := query{}
		q.Desc("name")
		if len(q.orderBy) != 1 {
			t.Fatal("err")
		}
	})
}

func Test_query_Page(test *testing.T) {
	test.Run("args:0", func(t *testing.T) {
		q := query{}
		q.Page(0)
		if q.page != 1 {
			t.Fatal("err")
		}
	})

	test.Run("success", func(t *testing.T) {
		q := query{}
		q.Page(2)
		if q.page != 2 {
			t.Fatal("err")
		}
	})
}

func Test_query_PageSize(test *testing.T) {
	test.Run("args:0", func(t *testing.T) {
		q := query{}
		q.PageSize(0)
		if q.pageSize != 10 {
			t.Fatal("err")
		}
	})

	test.Run("success", func(t *testing.T) {
		q := query{}
		q.PageSize(2)
		if q.pageSize != 2 {
			t.Fatal("err")
		}
	})
}

func Test_query_Where(test *testing.T) {
	test.Run("empty", func(t *testing.T) {
		q := query{}
		q.Where(bson.M{})
		if len(q.filter) != 0 {
			t.Fatal("err")
		}
	})

	test.Run("success", func(t *testing.T) {
		q := query{}
		q.Where(bson.M{"name": "test-name"})
		if len(q.filter) != 1 && q.filter["name"] != "test-name" {
			t.Fatal("err")
		}
	})
}

func Test_query_Find(test *testing.T) {
	repo, err := getClient()
	if err != nil {
		test.Fatal("err", err)
	}

	test.Run("all", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "find-test-001",
				Name: "find-test-name-001",
			},
			{
				ID:   "find-test-002",
				Name: "find-test-name-002",
			},
			{
				ID:   "find-test-003",
				Name: "find-test-name-003",
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		a := assert.New(t)
		var entries []testPerson
		if err = repo.Query().Find(&entries); err != nil {
			t.Fatal("err", err)
		}
		if len(entries) != 3 {
			t.Fatal("err", len(entries))
		}
		a.Equal(entries, addEntries)
	})

	test.Run("find.filter", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "find-filter-test-001",
				Name: "find-filter-test-name1",
			},
			{
				ID:   "find-filter-test-002",
				Name: "find-filter-test-name",
			},
			{
				ID:   "find-filter-test-003",
				Name: "find-filter-test-name",
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		var entries []testPerson
		if err = repo.Query().Where(bson.M{"name": addEntries[1].Name}).Find(&entries); err != nil {
			t.Fatal("err", err)
		}
		if len(entries) != 2 {
			t.Fatal("err", len(entries))
		}

		a := assert.New(t)
		compareEntries := []testPerson{
			addEntries[1],
			addEntries[2],
		}
		a.Equal(entries, compareEntries)
	})

	test.Run("order.filter.find", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "find-filter-test-001",
				Name: "find-filter-test-name1",
				Age:  20,
			},
			{
				ID:   "find-filter-test-002",
				Name: "find-filter-test-name",
				Age:  33,
			},
			{
				ID:   "find-filter-test-003",
				Name: "find-filter-test-name",
				Age:  15,
			},
			{
				ID:   "find-filter-test-004",
				Name: "find-filter-test-name",
				Age:  23,
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		var entries []testPerson
		if err = repo.Query().Desc("age").Where(bson.M{"name": addEntries[1].Name}).Find(&entries); err != nil {
			t.Fatal("err", err)
		}
		if len(entries) != 3 {
			t.Fatal("err", len(entries))
		}

		a := assert.New(t)
		compareEntries := []testPerson{
			addEntries[1],
			addEntries[3],
			addEntries[2],
		}
		a.Equal(entries, compareEntries)
	})

	test.Run("projection 1.filter.find", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "find-filter-test-001",
				Name: "find-filter-test-name1",
				Age:  20,
			},
			{
				ID:   "find-filter-test-002",
				Name: "find-filter-test-name",
				Age:  33,
			},
			{
				ID:   "find-filter-test-003",
				Name: "find-filter-test-name",
				Age:  15,
			},
			{
				ID:   "find-filter-test-004",
				Name: "find-filter-test-name",
				Age:  23,
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()

		expectedEntries := make([]testPerson, 0)
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
			expectedEntries = append(expectedEntries, testPerson{
				ID:  entry.ID,
				Age: entry.Age,
			})
		}

		var entries []testPerson
		if err = repo.Query().Fields(bson.M{"_id": 1, "age": 1}).Find(&entries); err != nil {
			t.Fatal("err", err)
		}

		a := assert.New(t)
		a.Equal(expectedEntries, entries)
	})
}

func Test_query_First(test *testing.T) {
	repo, err := getClient()
	if err != nil {
		test.Fatal("err", err)
	}

	test.Run("default.top1", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "first-test-001",
				Name: "first-test-name-001",
			},
			{
				ID:   "first-test-002",
				Name: "first-test-name-002",
			},
			{
				ID:   "first-test-003",
				Name: "first-test-name-003",
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		var entry testPerson
		if err = repo.Query().First(&entry); err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(addEntries[0], entry)
	})

	test.Run("desc.success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "first-test-001",
				Name: "first-test-name-001",
				Age:  1,
			},
			{
				ID:   "first-test-002",
				Name: "first-test-name-002",
				Age:  2,
			},
			{
				ID:   "first-test-003",
				Name: "first-test-name-003",
				Age:  3,
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		var entry testPerson
		if err = repo.Query().Desc("age").First(&entry); err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(addEntries[2], entry)
	})

	test.Run("filter.success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "first-test-001",
				Name: "first-test-name-001",
				Age:  15,
			},
			{
				ID:   "first-test-002",
				Name: "first-test-name-002",
				Age:  35,
			},
			{
				ID:   "first-test-003",
				Name: "first-test-name-003",
				Age:  26,
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		var entry testPerson
		if err = repo.Query().Where(bson.M{"age": bson.M{"$gt": 25}}).Asc("age").First(&entry); err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(addEntries[2], entry)
	})

	test.Run("projection 1.success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "first-test-001",
				Name: "first-test-name-001",
				Age:  15,
			},
			{
				ID:   "first-test-002",
				Name: "first-test-name-002",
				Age:  35,
			},
			{
				ID:   "first-test-003",
				Name: "first-test-name-003",
				Age:  26,
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		var entry testPerson
		if err = repo.Query().Fields(bson.M{"age": 1}).Where(bson.M{"_id": addEntries[1].ID}).First(&entry); err != nil {
			t.Fatal("err", err)
		}

		expectedEntry := testPerson{
			ID:  addEntries[1].ID,
			Age: addEntries[1].Age,
		}
		a := assert.New(t)
		a.Equal(expectedEntry, entry)
	})

	test.Run("projection 0.success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   "first-test-001",
				Name: "first-test-name-001",
				Age:  15,
			},
			{
				ID:   "first-test-002",
				Name: "first-test-name-002",
				Age:  35,
			},
			{
				ID:   "first-test-003",
				Name: "first-test-name-003",
				Age:  26,
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			entry := addEntries[index]
			_, err = collectionDb.InsertOne(ctx, &entry)
			if err != nil {
				t.Fatal("err", err)
			}
		}

		var entry testPerson
		if err = repo.Query().Fields(bson.M{"name": 0}).Where(bson.M{"_id": addEntries[1].ID}).First(&entry); err != nil {
			t.Fatal("err", err)
		}

		expectedEntry := testPerson{
			ID:  addEntries[1].ID,
			Age: addEntries[1].Age,
		}
		a := assert.New(t)
		a.Equal(expectedEntry, entry)
	})

	test.Run("primitive.ObjectID.success", func(t *testing.T) {
		addEntries := []objectEntry{
			{
				ID:   primitive.NewObjectID(),
				Desc: "first-desc",
			},
		}
		collectionDb := repo.database.Collection(addEntries[0].Table())
		ctx := context.Background()
		defer func() {
			for _, entry := range addEntries {
				_, _ = collectionDb.DeleteOne(ctx, bson.M{"_id": entry.ID})
			}
		}()
		for index := range addEntries {
			_, err = collectionDb.InsertOne(ctx, &addEntries[index])
			if err != nil {
				t.Fatal("err", err)
			}
		}

		var queryEntry objectEntry
		if err = repo.Query().Fields(bson.M{"name": 0}).Where(bson.M{"_id": addEntries[0].ID}).First(&queryEntry); err != nil {
			t.Fatal("err", err)
		}

		a := assert.New(t)
		a.Equal(addEntries[0], queryEntry)
	})
}
