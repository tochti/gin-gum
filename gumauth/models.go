package gumauth

import "time"

type (
	User struct {
		ID        int64     `db:"id" json:"id"`
		Username  string    `db:"username" json:"username"`
		FirstName string    `db:"first_name" json:"first_name"`
		LastName  string    `db:"last_name" json:"last_name"`
		Email     string    `db:"email" json:"email"`
		Password  string    `db:"password" json:"password"`
		IsActive  bool      `db:"is_active" json:"is_active"`
		LastLogin time.Time `db:"last_login" json:"last_login"`
	}
)
