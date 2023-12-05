package mysqlex

import (
	"testing"

	"gorm.io/gorm"
)

type UserValue struct {
	ID    int `json:"id" gorm:"column:id;primaryKey;autoincrement;type:INT;comment:数据标识"`
	Value int `json:"value" gorm:"column:value;type:int"`
}

func (m UserValue) TableName() string {
	return "user_value"
}

func Text_UPdate(t *testing.T) {
	db, _ := getDb()
	err := db.Transaction(func(tx *gorm.DB) error {
		for i := 1; i <= 5; i++ {

		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
