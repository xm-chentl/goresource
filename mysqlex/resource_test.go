package mysqlex

import (
	"testing"

	"github.com/xm-chentl/goresource"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testConnectionStr = "root:987654321@tcp(47.98.248.82:9902)/my_test"
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
