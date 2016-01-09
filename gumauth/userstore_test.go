package gumauth

import (
	"testing"
)

func Test_NewSha512Password_OK(t *testing.T) {
	tokens := []string{}
	for x := 0; x < 10; x++ {
		token := NewSha512Password(string(x))
		if !ExistsToken(tokens, token) {
			t.Fatal("Expect every token to be unique", token)
		}
		tokens = append(tokens, token)
	}
}

func Test_StoreUser(t *testing.T) {
	user := User{
		ID:       1,
		Username: "joe",
		Password: NewSha512Password("cameron"),
	}
	su := NewStoreUser(user)

	if !su.ValidPassword("cameron") {
		t.Fatal("Expect password to be correct")
	}

	if su.ValidPassword("none") {
		t.Fatal("Expect password to be uncorrect")
	}

	if su.ID() != "1" {
		t.Fatal("Expect id to be 1")
	}

	if su.Password() != NewSha512Password("cameron") {
		t.Fatal("Expect password to be cameron")
	}
}

func Test_SQLUserStore(t *testing.T) {
	db := initTestDB(t)

	users := []*User{
		{
			Username: "foo",
			Password: "bar",
		},
		{
			Username: "anonym",
			Password: "123",
		},
	}

	err := db.Insert(users[0], users[1])
	if err != nil {
		t.Fatal(err)
	}

	store := SQLUserStore{db.Db}

	u, err := store.FindUser("foo")
	if err != nil {
		t.Fatal(err)
	}

	if u.ID() != "1" {
		t.Fatal("Expect %v was %v", 1, u.ID())
	}
}
