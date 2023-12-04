package mysqlex

import (
	"testing"

	"github.com/xm-chentl/goresource"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testConnectionStr = "root:123456@tcp(localhost:3306)/test_db"
var testDb *gorm.DB

func Test_Connection(t *testing.T) {
	var err error
	testDb, err = gorm.Open(mysql.Open(testConnectionStr), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
}

func getDb() (*gorm.DB, goresource.IRepository) {
	if testDb == nil {
		testDb, _ = gorm.Open(mysql.Open(testConnectionStr), &gorm.Config{})
	}

	return testDb, &repository{
		db: testDb,
	}
}
