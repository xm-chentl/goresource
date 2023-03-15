package mysqlex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPerson struct {
	Name string `gorm:"column:name"`
	ID   int64  `gorm:"column:id;primaryKey"`
	Age  int8   `gorm:column:age"`
}

func (m TestPerson) GetID() interface{} {
	return m.ID
}

func (m *TestPerson) SetID(v interface{}) {
	if v != nil {
		m.ID = v.(int64)
	}
}

func (m TestPerson) Table() string {
	return "test_person"
}

func (m TestPerson) TableName() string {
	return m.Table()
}

type TestIdentity struct {
	ID   uint64 `gorm:"column:id;AUTO_INCREMENT;primaryKey"`
	Name string `gorm:column:"name"`
}

func (m TestIdentity) GetID() interface{} {
	return m.ID
}

func (m *TestIdentity) SetID(v interface{}) {
	if v != nil {
		m.ID = v.(uint64)
	}
}

func (m TestIdentity) Table() string {
	return "test_identity"
}

func (m TestIdentity) TableName() string {
	return m.Table()
}

func Test_Create(test *testing.T) {
	db, repoDb := getDb()
	test.Run("success", func(t *testing.T) {
		var err error
		entry := TestPerson{
			ID:   123456,
			Name: "tenglong.chen",
			Age:  36,
		}
		if err = repoDb.Create(&entry); err != nil {
			t.Fatal(err)
		}

		newEntry := TestPerson{}
		if err = db.Model(&newEntry).Where("id = ?", entry.ID).First(&newEntry).Error; err != nil {
			t.Fatal(err)
		}
		if newEntry.ID != entry.ID {
			t.Fatal("err")
		}
		_ = db.Model(entry).Delete(entry).Error
	})

	test.Run("identity success", func(t *testing.T) {
		var err error
		entry := TestIdentity{
			Name: "tenglong.chen",
		}
		if err = repoDb.Create(&entry); err != nil {
			t.Fatal(err)
		}
		if entry.ID == 0 {
			t.Fatal("identity id is empty")
		}

		newEntry := TestIdentity{}
		if err = db.Model(&newEntry).Where("id = ?", entry.ID).First(&newEntry).Error; err != nil {
			t.Fatal(err)
		}
		_ = db.Model(entry).Delete(entry).Error
	})
}

func Test_Delete(test *testing.T) {
	db, repoDb := getDb()
	test.Run("success", func(t *testing.T) {
		var err error
		entries := []TestPerson{
			{
				ID:   123,
				Name: "tenglong.chen",
				Age:  36,
			},
			{
				ID:   456,
				Name: "tenglong.chen2",
				Age:  37,
			},
		}
		for index := range entries {
			if err = db.Model(&entries[index]).Create(&entries[index]).Error; err != nil {
				t.Fatal(err)
			}
		}
		defer func() {
			for index := range entries {
				if err = db.Delete(&entries[index]).Error; err != nil {
					t.Fatal(err)
				}
			}
		}()

		delEntry := TestPerson{
			ID: 123,
		}
		if err = repoDb.Delete(&delEntry); err != nil {
			t.Fatal(err)
		}
	})
}

func Test_Update(test *testing.T) {
	db, repoDb := getDb()
	test.Run("success a single field", func(t *testing.T) {
		var err error
		entries := []TestPerson{
			{
				ID:   789,
				Name: "tenglong.chen",
				Age:  36,
			},
			{
				ID:   789123,
				Name: "tenglong.chen2",
				Age:  37,
			},
		}
		for index := range entries {
			if err = db.Create(&entries[index]).Error; err != nil {
				t.Fatal(err)
			}
		}
		defer func() {
			for index := range entries {
				if err = db.Delete(&entries[index]).Error; err != nil {
					t.Fatal(err)
				}
			}
		}()

		updateEntry := TestPerson{
			ID:   789123,
			Name: "minghao.chen",
			Age:  11,
		}
		if err = repoDb.Update(&updateEntry); err != nil {
			t.Fatal(err)
		}

		updatedEntry := TestPerson{}
		if err = db.Model(&updatedEntry).Where("id = ?", updateEntry.ID).First(&updatedEntry).Error; err != nil {
			t.Fatal(err)
		}

		a := assert.New(t)
		a.Equal(updateEntry, updatedEntry)
	})

	test.Run("success omit field", func(t *testing.T) {
		var err error
		entries := []TestPerson{
			{
				ID:   789,
				Name: "tenglong.chen",
				Age:  36,
			},
			{
				ID:   789123,
				Name: "tenglong.chen2",
				Age:  37,
			},
		}
		for index := range entries {
			if err = db.Create(&entries[index]).Error; err != nil {
				t.Fatal(err)
			}
		}
		defer func() {
			for index := range entries {
				if err = db.Delete(&entries[index]).Error; err != nil {
					t.Fatal(err)
				}
			}
		}()

		updateEntry := TestPerson{
			ID:   789123,
			Name: "minghao.chen",
			Age:  11,
		}
		if err = repoDb.Update(&updateEntry, SaveOptionByOmit{
			Fields: []string{"age"},
		}); err != nil {
			t.Fatal(err)
		}

		updatedEntry := TestPerson{}
		if err = db.Model(&updatedEntry).Where("id = ?", updateEntry.ID).First(&updatedEntry).Error; err != nil {
			t.Fatal(err)
		}

		a := assert.New(t)
		a.Equal(TestPerson{
			ID:   updateEntry.ID,
			Name: updateEntry.Name,
			Age:  37,
		}, updatedEntry)
	})
}
