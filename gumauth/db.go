package gumauth

import (
	"database/sql"
	"fmt"

	"gopkg.in/gorp.v1"
)

var (
	UserTable = "users"
)

func CreateTables(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}

	q := `
	CREATE TABLE IF NOT EXISTS %v (
		id BIGINT(20) PRIMARY KEY AUTO_INCREMENT,
		username VARCHAR(255) NOT NULL,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		email VARCHAR(500),
		password VARCHAR(255),
		is_active BOOL,
		last_login DATETIME
	)
	`

	query := fmt.Sprintf(q, UserTable)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func AddTables(db *gorp.DbMap) {
	db.AddTableWithName(User{}, UserTable).SetKeys(true, "id")
}
