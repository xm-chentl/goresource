package postgres

import (
	"context"
	"testing"

	"github.com/xm-chentl/goresource/postgres/grammar"
	"github.com/xm-chentl/goresource/postgres/metadata"
	"github.com/xm-chentl/goresource/errs"

	"github.com/stretchr/testify/assert"
)

var (
	createTestPesonSql = `drop table if exists test_person;CREATE TABLE test_person (
		id int8 NOT NULL,
		"name" varchar NULL,
		age int2 NULL,
		CONSTRAINT test_person_pk PRIMARY KEY (id)
	);`
	// 放至最后一个test文件
	dropTestPersonSql = `drop table if exists test_person;`
)

type testPerson struct {
	ID   int64  `postgres:"id" pk:""`
	Name string `postgres:"name"`
	Age  int16  `postgres:"age"`
}

func (t testPerson) GetID() interface{} {
	return t.ID
}

func (t *testPerson) SetID(v interface{}) {
	if v == nil {
		return
	}
	t.ID = v.(int64)
}

func (t testPerson) Table() string {
	return "test_person"
}

func Test_repository_Create(t *testing.T) {
	repo, err := getRepo()
	if err != nil {
		t.Fatal("err", err)
	}
	addEntries := []testPerson{
		{
			ID:   10001,
			Name: "create_test-001",
			Age:  21,
		},
		{
			ID:   10002,
			Name: "create_test_002",
			Age:  11,
		},
	}

	conn, err := repo.pool.getConn()
	if err != nil {
		t.Fatal("err", err)
	}
	ctx := context.Background()
	defer func() {
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Delete(metadata.Get(&entry), &entry, "id = $1", entry.GetID())
			_, err = conn.Exec(ctx, sql, args...)
			if err != nil {
				t.Fatal("err", err, sql, args)
			}
		}
		conn.Release()
	}()
	for index := range addEntries {
		entry := addEntries[index]
		if err = repo.Create(&entry); err != nil {
			t.Fatal("err", err)
		}
	}

	sql, args := grammar.Select(metadata.Get(&addEntries[0]), []string{})
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		t.Fatal("err", err, sql)
	}
	entries := make([]testPerson, 0)
	for rows.Next() {
		entry := testPerson{}
		if err = rows.Scan(&entry.ID, &entry.Name, &entry.Age); err != nil {
			t.Fatal("err", err)
		}
		entries = append(entries, entry)
	}
	a := assert.New(t)
	a.Equal(addEntries, entries)
}

func Test_repository_Delete(test *testing.T) {
	repo, err := getRepo()
	if err != nil {
		test.Fatal("err", err)
	}

	test.Run(errs.DeleteFullNotAllowed.Error(), func(t *testing.T) {
		deleteEntry := testPerson{
			ID:   0,
			Name: "delete_test_001",
			Age:  1,
		}
		err := repo.Delete(&deleteEntry)
		a := assert.New(t)
		a.Equal(errs.DeleteFullNotAllowed, err)
	})

	test.Run("one delete by entry", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10331,
				Name: "delete_test-001",
				Age:  21,
			},
			{
				ID:   10332,
				Name: "delete_test_002",
				Age:  11,
			},
		}
		conn, err := repo.pool.getConn()
		if err != nil {
			t.Fatal("err", err)
		}

		ctx := context.Background()
		defer func() {
			for index := range addEntries {
				entry := addEntries[index]
				sql, args := grammar.Delete(metadata.Get(&entry), &entry, "id = $1", entry.GetID())
				_, err = conn.Exec(ctx, sql, args...)
				if err != nil {
					t.Fatal("err", err, sql, args)
				}
			}
			conn.Release()
		}()
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
		}
		if err = repo.Delete(&addEntries[0]); err != nil {
			t.Fatal("err", err)
		}

		sql, args := grammar.Select(metadata.Get(&addEntries[0]), []string{}, "id = $1", addEntries[0].ID)
		rows, err := conn.Query(ctx, sql, args...)
		if err != nil {
			t.Fatal("err", err, sql)
		}

		entries := make([]testPerson, 0)
		for rows.Next() {
			entry := testPerson{}
			if err = rows.Scan(&entry.ID, &entry.Name, &entry.Age); err != nil {
				t.Fatal("err", err)
			}
			entries = append(entries, entry)
		}
		a := assert.New(t)
		a.Equal(0, len(entries))
	})

	test.Run("many delete by entry", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10341,
				Name: "delete_test-001",
				Age:  21,
			},
			{
				ID:   10342,
				Name: "delete_test_002",
				Age:  11,
			},
			{
				ID:   10343,
				Name: "delete_test_003",
				Age:  11,
			},
		}
		conn, err := repo.pool.getConn()
		if err != nil {
			t.Fatal("err", err)
		}

		ctx := context.Background()
		defer func() {
			for index := range addEntries {
				entry := addEntries[index]
				sql, args := grammar.Delete(metadata.Get(&entry), &entry, "id = $1", entry.GetID())
				_, err = conn.Exec(ctx, sql, args...)
				if err != nil {
					t.Fatal("err", err, sql, args)
				}
			}
			conn.Release()
		}()
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
		}
		if err = repo.Delete(&addEntries[0], "age = 11"); err != nil {
			t.Fatal("err", err)
		}

		sql, args := grammar.Select(metadata.Get(&addEntries[0]), []string{}, "age = 11")
		rows, err := conn.Query(ctx, sql, args...)
		if err != nil {
			t.Fatal("err", err, sql)
		}

		entries := make([]testPerson, 0)
		for rows.Next() {
			entry := testPerson{}
			if err = rows.Scan(&entry.ID, &entry.Name, &entry.Age); err != nil {
				t.Fatal("err", err)
			}
			entries = append(entries, entry)
		}
		a := assert.New(t)
		a.Equal(0, len(entries))
	})
}

func Test_repository_Update(test *testing.T) {
	repo, err := getRepo()
	if err != nil {
		test.Fatal("err", err)
	}

	// 暂时逻辑，不存在全量更新
	// test.Run(errs.UpdateFullNotAllowed.Error(), func(t *testing.T) {
	// 	updateEntry := testPerson{
	// 		ID:   0,
	// 		Name: "update-set-name",
	// 		Age:  11,
	// 	}
	// 	err = repo.Update(&updateEntry, []string{"name"})
	// 	a := assert.New(t)
	// 	a.Equal(errs.UpdateFullNotAllowed, err)
	// })

	test.Run("one update.set.fields by entry", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10011,
				Name: "update_test-001",
				Age:  21,
			},
			{
				ID:   10012,
				Name: "update_test_002",
				Age:  11,
			},
		}
		conn, err := repo.pool.getConn()
		if err != nil {
			t.Fatal("err", err)
		}

		ctx := context.Background()
		defer func() {
			for index := range addEntries {
				entry := addEntries[index]
				sql, args := grammar.Delete(metadata.Get(&entry), &entry, "id = $1", entry.GetID())
				_, err = conn.Exec(ctx, sql, args...)
				if err != nil {
					t.Fatal("err", err, sql, args)
				}
			}
			conn.Release()
		}()
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
		}

		updateEntry := testPerson{
			ID:   addEntries[0].ID,
			Name: "update-set-name",
			Age:  addEntries[0].Age,
		}
		if err = repo.Update(&updateEntry, []string{"name"}, "id = $2", updateEntry.ID); err != nil {
			t.Fatal("err", err)
		}

		sql, args := grammar.Select(metadata.Get(&updateEntry), []string{}, "id = $1", updateEntry.ID)
		rows, err := conn.Query(ctx, sql, args...)
		if err != nil {
			t.Fatal("err", err, sql)
		}

		entries := make([]testPerson, 0)
		for rows.Next() {
			entry := testPerson{}
			if err = rows.Scan(&entry.ID, &entry.Name, &entry.Age); err != nil {
				t.Fatal("err", err)
			}
			entries = append(entries, entry)
		}
		a := assert.New(t)
		a.Equal(1, len(entries))
		a.Equal(entries[0], updateEntry)
	})

	test.Run("many entry update.set.fields by filter", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10031,
				Name: "update_test_filter_001",
				Age:  21,
			},
			{
				ID:   10032,
				Name: "update_test_filter_002",
				Age:  11,
			},
			{
				ID:   10033,
				Name: "update_test_filter_003",
				Age:  11,
			},
		}
		conn, err := repo.pool.getConn()
		if err != nil {
			t.Fatal("err", err)
		}

		ctx := context.Background()
		defer func() {
			for index := range addEntries {
				entry := addEntries[index]
				sql, args := grammar.Delete(metadata.Get(&entry), &entry, "id = $1", entry.GetID())
				_, err = conn.Exec(ctx, sql, args...)
				if err != nil {
					t.Fatal("err", err, sql, args)
				}
			}
			conn.Release()
		}()
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
		}

		updateEntry := testPerson{
			ID:   addEntries[0].ID,
			Name: "update-set-name",
			Age:  11,
		}
		if err = repo.Update(&updateEntry, []string{"name"}, "age = $2", updateEntry.Age); err != nil {
			t.Fatal("err", err)
		}

		sql, args := grammar.Select(metadata.Get(&updateEntry), []string{}, "age = 11")
		rows, err := conn.Query(ctx, sql, args...)
		if err != nil {
			t.Fatal("err", err, sql, args)
		}

		entries := make([]testPerson, 0)
		for rows.Next() {
			entry := testPerson{}
			if err = rows.Scan(&entry.ID, &entry.Name, &entry.Age); err != nil {
				t.Fatal("err", err)
			}
			entries = append(entries, entry)
		}
		a := assert.New(t)
		a.Equal(2, len(entries))
		a.Equal(entries, []testPerson{
			{
				ID:   10032,
				Name: "update-set-name",
				Age:  11,
			},
			{
				ID:   10033,
				Name: "update-set-name",
				Age:  11,
			},
		})
	})

	test.Run("many entry update.set.fields by many filter", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10041,
				Name: "update_test_filter_001",
				Age:  21,
			},
			{
				ID:   10042,
				Name: "update_test_filter_002",
				Age:  11,
			},
			{
				ID:   10043,
				Name: "update_test_filter_003",
				Age:  11,
			},
			{
				ID:   10044,
				Name: "update_test_filter",
				Age:  31,
			},
			{
				ID:   10045,
				Name: "update_test_filter",
				Age:  31,
			},
		}
		conn, err := repo.pool.getConn()
		if err != nil {
			t.Fatal("err", err)
		}

		ctx := context.Background()
		defer func() {
			for index := range addEntries {
				entry := addEntries[index]
				sql, args := grammar.Delete(metadata.Get(&entry), &entry, "id = $1", entry.GetID())
				_, err = conn.Exec(ctx, sql, args...)
				if err != nil {
					t.Fatal("err", err, sql, args)
				}
			}
			conn.Release()
		}()
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
		}

		updateEntry := testPerson{
			ID:   addEntries[0].ID,
			Name: "update-set-name",
			Age:  31,
		}
		if err = repo.Update(&updateEntry, []string{"name"}, "name = $2 AND age = $3", addEntries[3].Name, updateEntry.Age); err != nil {
			t.Fatal("err", err)
		}

		sql, args := grammar.Select(metadata.Get(&updateEntry), []string{}, "age = $1 AND name = $2", 31, updateEntry.Name)
		rows, err := conn.Query(ctx, sql, args...)
		if err != nil {
			t.Fatal("err", err, sql, args)
		}

		entries := make([]testPerson, 0)
		for rows.Next() {
			entry := testPerson{}
			if err = rows.Scan(&entry.ID, &entry.Name, &entry.Age); err != nil {
				t.Fatal("err", err)
			}
			entries = append(entries, entry)
		}
		a := assert.New(t)
		a.Equal(2, len(entries))
		a.Equal(entries, []testPerson{
			{
				ID:   10044,
				Name: "update-set-name",
				Age:  31,
			},
			{
				ID:   10045,
				Name: "update-set-name",
				Age:  31,
			},
		})
	})
}

func Test_end(t *testing.T) {
	repo, err := getRepo()
	if err != nil {
		t.Fatal("err", err)
	}

	_, err = repo.pool.pgxPool.Exec(repo.ctx, dropTestPersonSql)
	if err != nil {
		t.Fatal("err", err)
	}

	_, err = repo.pool.pgxPool.Exec(repo.ctx, dropTestTimeSql)
	if err != nil {
		t.Fatal("err", err)
	}
}
