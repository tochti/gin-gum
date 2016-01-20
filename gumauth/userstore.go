package gumauth

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"strconv"
)

type (
	StoreUser interface {
		ID() string
		Password() string
		ValidPassword(string) bool
	}

	UserStore interface {
		FindUser(string) (StoreUser, error)
	}

	storeUser struct {
		user User
	}

	SQLUserStore struct {
		DB *sql.DB
	}
)

func NewStoreUser(u User) StoreUser {
	return storeUser{
		user: u,
	}
}

func (u storeUser) ID() string {
	return strconv.FormatInt(u.user.ID, 10)
}

func (u storeUser) Password() string {
	return u.user.Password
}

// Check if given password correct. Passed password get convert to sha512 hash
func (u storeUser) ValidPassword(p string) bool {
	return NewSha512Password(p) == u.Password()
}

func (s SQLUserStore) FindUser(name string) (StoreUser, error) {
	u := User{}
	q := fmt.Sprintf(`SELECT id,username,password 
			  FROM %v
			  WHERE username=?`, UserTable)
	err := s.DB.QueryRow(q, name).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}

	return NewStoreUser(u), nil
}

func NewSha512Password(pass string) string {
	hash := sha512.New()
	tmp := hash.Sum([]byte(pass))
	passHash := fmt.Sprintf("%x", tmp)
	return passHash
}
