package gumspecs

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQL(t *testing.T) {
	os.Clearenv()
	os.Setenv("T_MYSQL_USER", "foo")
	os.Setenv("T_MYSQL_PASSWORD", "bar")
	os.Setenv("T_MYSQL_HOST", "h")
	os.Setenv("T_MYSQL_PORT", "p")
	os.Setenv("T_MYSQL_DB_NAME", "n")
	os.Setenv("T_MYSQL_LOCATION", "l")

	AppName = "t"
	mysql := ReadMySQL()
	if mysql == nil {
		t.Fatal("Expect *MySQL Type was nil")
	}

	_, err := mysql.DB()
	if err != nil {
		t.Fatal("Expect nil was %v", err)
	}

	mysqlStr := "foo:bar@tcp(h:p)/n?loc=l&parseTime=true"
	if mysql.String() != mysqlStr {
		t.Fatalf("Expect %v was %v", mysqlStr, mysql.String())
	}
}
