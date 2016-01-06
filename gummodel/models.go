package gummodel

import "time"

type (
	User struct {
		ID        int64     `db:"id"`
		Username  string    `db:"username"`
		FirstName string    `db:"first_name"`
		LastName  string    `db:"last_name"`
		Email     string    `db:"email"`
		Password  string    `db:"password"`
		IsActive  bool      `db:"is_active"`
		LastLogin time.Time `db:"last_login"`
	}
)
