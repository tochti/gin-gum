package gumauth

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tochti/gin-gum/gumspecs"
	"gopkg.in/gorp.v1"
)

const (
	TestDatabase = "testing"
)

func setenvMySQL() {
	os.Clearenv()
	os.Setenv("MYSQL_USER", "tochti")
	os.Setenv("MYSQL_PASSWORD", "123")
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB_NAME", TestDatabase)
}

func initTestDB(t *testing.T) *gorp.DbMap {
	setenvMySQL()

	mysql := gumspecs.ReadMySQL()

	sqlDB, err := mysql.DB()
	if err != nil {
		t.Fatal(err)
	}

	db := &gorp.DbMap{
		Db:      sqlDB,
		Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"},
	}

	db.AddTableWithName(User{}, UserTable).SetKeys(true, "id")

	err = db.DropTablesIfExists()
	if err != nil {
		t.Fatal(err)
	}

	err = db.CreateTables()
	if err != nil {
		t.Fatal(err)
	}

	return db
}
