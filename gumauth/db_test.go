package gumauth

import (
	"testing"

	"gopkg.in/gorp.v1"

	"github.com/tochti/gin-gum/gumspecs"
)

func TestCreateTable(t *testing.T) {
	setenvMySQL()
	mysql := gumspecs.ReadMySQL()

	db, err := mysql.DB()
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTables(db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddTablesGorp(t *testing.T) {
	setenvMySQL()
	mysql := gumspecs.ReadMySQL()

	db, err := mysql.DB()
	if err != nil {
		t.Fatal(err)
	}

	dbMap := &gorp.DbMap{
		Db: db,
		Dialect: gorp.MySQLDialect{
			"UTF8",
			"InnonDB",
		},
	}

	AddTables(dbMap)
}
