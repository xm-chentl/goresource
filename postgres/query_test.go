package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/xm-chentl/goresource/postgres/grammar"
	"github.com/xm-chentl/goresource/postgres/metadata"

	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
)

var (
	createTimeStructSql = `drop table if exists test_time;CREATE TABLE test_time (
		id int8 NULL,
		"time" timestamptz NULL,
		value varchar NULL,
		CONSTRAINT test_time_pk PRIMARY KEY (id)
	);`
	// 放至最后一个test文件
	dropTestTimeSql = `drop table if exists test_person;`
)

type testTimeEntry struct {
	ID    int64     `postgres:"id" pk:""`
	Time  time.Time `postgres:"time"`
	Value string    `postgres:"value"`
}

func (t testTimeEntry) GetID() interface{} {
	return t.ID
}

func (t *testTimeEntry) SetID(v interface{}) {
	if vv, ok := v.(int64); ok {
		t.ID = vv
	}
}

func (t testTimeEntry) Table() string {
	return "test_time"
}

type tagPersonRes struct {
	ID   int64  `postgres:"id"`
	Name string `postgres:"name"`
}
type noTagPersonRes struct {
	ID   int64
	Name string
}

func Test_DataType(t *testing.T) {
	var targetValue interface{}
	value := &pgtype.Numeric{}
	value.Set(8)
	targetValue = value
	v, ok := targetValue.(pgtype.Value)
	if !ok {
		t.Fatal("err")
	}

	var tv int64
	if err := v.AssignTo(&tv); err != nil {
		t.Fatal(err)
	}
	if tv != 8 {
		t.Fatal("err", tv)
	}
}

func Test_query_Exec(test *testing.T) {
	repo, err := getRepo()
	if err != nil {
		test.Fatal("err", err)
	}

	// 有带tag用做映射
	test.Run("success.struct.carry tag", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10841,
				Name: "query_exec_001",
				Age:  21,
			},
			{
				ID:   10842,
				Name: "query_exec_002",
				Age:  11,
			},
			{
				ID:   10843,
				Name: "query_exec_003",
				Age:  11,
			},
			{
				ID:   10844,
				Name: "query_exec_004",
				Age:  31,
			},
			{
				ID:   10845,
				Name: "query_exec_005",
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
		expectRes := make([]tagPersonRes, 0)
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
			expectRes = append(expectRes, tagPersonRes{
				ID:   entry.ID,
				Name: entry.Name,
			})
		}

		var res []tagPersonRes
		if err = repo.Query().Exec(&res, "SELECT id, name FROM test_person"); err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(expectRes, res)
	})

	test.Run("success.struct.carry tag select fields", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10841,
				Name: "query_exec_001",
				Age:  21,
			},
			{
				ID:   10842,
				Name: "query_exec_002",
				Age:  11,
			},
			{
				ID:   10843,
				Name: "query_exec_003",
				Age:  11,
			},
			{
				ID:   10844,
				Name: "query_exec_004",
				Age:  31,
			},
			{
				ID:   10845,
				Name: "query_exec_005",
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
		expectRes := make([]tagPersonRes, 0)
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
			expectRes = append(expectRes, tagPersonRes{
				ID:   entry.ID,
				Name: entry.Name,
			})
		}

		var res []tagPersonRes
		if err = repo.Query().Exec(&res, "SELECT id, age, name FROM test_person"); err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(expectRes, res)
	})

	// 没带tag用字段名映射
	test.Run("success.struct.carry tag no", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10851,
				Name: "query_exec_struct_001",
				Age:  21,
			},
			{
				ID:   10852,
				Name: "query_exec_struct_002",
				Age:  11,
			},
			{
				ID:   10853,
				Name: "query_exec_struct_003",
				Age:  11,
			},
			{
				ID:   10854,
				Name: "query_exec_struct_004",
				Age:  31,
			},
			{
				ID:   10855,
				Name: "query_exec_struct_005",
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
		expectRes := make([]noTagPersonRes, 0)
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
			expectRes = append(expectRes, noTagPersonRes{
				ID:   entry.ID,
				Name: entry.Name,
			})
		}

		var res []noTagPersonRes
		if err = repo.Query().Exec(&res, "SELECT id, name FROM test_person"); err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(expectRes, res)
	})

	test.Run("success.filter.struct.carry tag no", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10861,
				Name: "query_exec_struct_001",
				Age:  21,
			},
			{
				ID:   10862,
				Name: "query_exec_struct_002",
				Age:  11,
			},
			{
				ID:   10863,
				Name: "query_exec_struct_003",
				Age:  11,
			},
			{
				ID:   10864,
				Name: "query_exec_struct_004",
				Age:  31,
			},
			{
				ID:   10865,
				Name: "query_exec_struct_005",
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
		expectRes := make([]noTagPersonRes, 0)
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
			if entry.Age == 11 {
				expectRes = append(expectRes, noTagPersonRes{
					ID:   entry.ID,
					Name: entry.Name,
				})
			}
		}

		var res []noTagPersonRes
		if err = repo.Query().Exec(&res, "SELECT id, name FROM test_person WHERE age = $1", 11); err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(2, len(res))
		a.Equal(expectRes, res)
	})
}

func Test_query_First(test *testing.T) {
	repo, err := getRepo()
	if err != nil {
		test.Fatal("err", err)
	}

	test.Run("success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10841,
				Name: "query_first_001",
				Age:  21,
			},
			{
				ID:   10842,
				Name: "query_first_002",
				Age:  11,
			},
			{
				ID:   10843,
				Name: "query_first_003",
				Age:  11,
			},
			{
				ID:   10844,
				Name: "query_first_004",
				Age:  31,
			},
			{
				ID:   10845,
				Name: "query_first_005",
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

		var entry testPerson
		if err = repo.Query().Where("id = 10845").First(&entry); err != nil {
			t.Fatal("err", err)
		}

		a := assert.New(t)
		a.Equal(addEntries[4], entry)
	})

	test.Run("success.fields", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10841,
				Name: "query_first_001",
				Age:  21,
			},
			{
				ID:   10842,
				Name: "query_first_002",
				Age:  11,
			},
			{
				ID:   10843,
				Name: "query_first_003",
				Age:  11,
			},
			{
				ID:   10844,
				Name: "query_first_004",
				Age:  31,
			},
			{
				ID:   10845,
				Name: "query_first_005",
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

		var entry testPerson
		if err = repo.Query().Fields("id", "age").Where("id = $1", addEntries[4].ID).First(&entry); err != nil {
			t.Fatal("err", err)
		}

		expectedEntry := testPerson{
			ID:  addEntries[4].ID,
			Age: addEntries[4].Age,
		}
		a := assert.New(t)
		a.Equal(expectedEntry, entry)
	})
}

func Test_query_Count(test *testing.T) {
	repo, err := getRepo()
	if err != nil {
		test.Fatal("err", err)
	}

	test.Run("success", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10941,
				Name: "query_count_001",
				Age:  21,
			},
			{
				ID:   10942,
				Name: "query_count_002",
				Age:  11,
			},
			{
				ID:   10943,
				Name: "query_count_003",
				Age:  11,
			},
			{
				ID:   10944,
				Name: "query_count_004",
				Age:  31,
			},
			{
				ID:   10945,
				Name: "query_count_005",
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

		count, err := repo.Query().Count(&testPerson{})
		if err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(int64(5), count)
	})

	test.Run("success.filter", func(t *testing.T) {
		addEntries := []testPerson{
			{
				ID:   10941,
				Name: "query_count_001",
				Age:  21,
			},
			{
				ID:   10942,
				Name: "query_count_002",
				Age:  11,
			},
			{
				ID:   10943,
				Name: "query_count_003",
				Age:  11,
			},
			{
				ID:   10944,
				Name: "query_count_004",
				Age:  31,
			},
			{
				ID:   10945,
				Name: "query_count_005",
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

		var expectedCount int64
		for index := range addEntries {
			entry := addEntries[index]
			sql, args := grammar.Insert(metadata.Get(&entry), &entry)
			if _, err = conn.Exec(ctx, sql, args...); err != nil {
				t.Fatal("err", err)
			}
			if entry.Age == addEntries[4].Age {
				expectedCount++
			}
		}

		count, err := repo.Query().Where("age = $1", addEntries[4].Age).Count(&testPerson{})
		if err != nil {
			t.Fatal("err", err)
		}
		a := assert.New(t)
		a.Equal(expectedCount, count)
	})
}

// todo: 未使用
func mockFunc(
	t *testing.T,
	repo repository,
	mockData []testPerson,
	rangeDataFunc func(index int, entry testPerson),
	cbFunc func(t *testing.T),
) {
	conn, err := repo.pool.getConn()
	if err != nil {
		t.Fatal("err", err)
	}

	ctx := context.Background()
	defer func() {
		for index := range mockData {
			entry := mockData[index]
			if rangeDataFunc != nil {
				rangeDataFunc(index, entry)
			}
			sql, args := grammar.Delete(metadata.Get(&entry), &entry, "id = $1", entry.GetID())
			_, err = conn.Exec(ctx, sql, args...)
			if err != nil {
				t.Fatal("err", err, sql, args)
			}

		}
		conn.Release()
	}()
	for index := range mockData {
		entry := mockData[index]
		sql, args := grammar.Insert(metadata.Get(&entry), &entry)
		if _, err = conn.Exec(ctx, sql, args...); err != nil {
			t.Fatal("err", err)
		}
	}
	if cbFunc != nil {
		cbFunc(t)
	}
}
