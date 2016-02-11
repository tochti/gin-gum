package gumauth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumtest"
)

func Test_SQLNewUser(t *testing.T) {
	db := initTestDB(t)
	user := newTestUser()

	err := SQLNewUser(db.Db, &user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 1 {
		t.Fatal("Expect user id 1 was", user.ID)
	}

}

func Test_SQLNewUserAlreadyExists(t *testing.T) {
	db := initTestDB(t)
	user := fillTestDB(t, db)

	err := SQLNewUser(db.Db, &user)
	if err != UserExistsErr {
		t.Fatal("Expect %v was %v", UserExistsErr, err)
	}

}

func Test_POST_User_OK(t *testing.T) {
	db := initTestDB(t)
	user := newTestUser()

	r := gin.New()
	r.POST("/", CreateUserSQL(db.Db))

	body := `
	{
		"username": "devilXX",
		"password": "123",
		"first_name": "Dare",
		"last_name": "Devil",
		"email": "devil@hell.de"
	}
	`

	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", body)
	if resp.Code != http.StatusCreated {
		t.Fatalf("Expect %v was %v", http.StatusCreated, resp.Code)
	}

	respUser := User{}
	err := json.Unmarshal(resp.Body.Bytes(), &respUser)
	if err != nil {
		t.Fatal(err)
	}

	if user.Username != respUser.Username ||
		user.Email != respUser.Email ||
		user.FirstName != respUser.FirstName ||
		user.LastName != respUser.LastName {
		t.Fatalf("Expect %v was %v", user, respUser)
	}

}
func Test_SignIn_OK(t *testing.T) {
	user := User{
		ID:       1,
		Username: "ladykiller_XX",
		Password: NewSha512Password("123"),
	}

	userStore := TestUserStore{user}
	sessStore := TestSessionStore{}

	r := gin.New()
	h := SignIn(sessStore, userStore)
	r.GET("/:name/:password", h)

	name := base64.StdEncoding.EncodeToString([]byte(user.Username))
	pass := base64.StdEncoding.EncodeToString([]byte("123"))

	url := fmt.Sprintf("/%v/%v", name, pass)
	resp := gumtest.NewRouter(r).ServeHTTP("GET", url, "")

	respSess := Session{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respSess); err != nil {
		t.Fatal(err)
	}

	if resp.Code != http.StatusAccepted {
		t.Fatalf("Expect %v was %v", resp.Code, http.StatusAccepted)
	}

	if respSess.Token != "1234" {
		t.Fatalf("Expect %v was %v", "1234", respSess.Token)
	}
	if respSess.UserID != "1" {
		t.Fatalf("Expect %v was %v", "1", respSess.UserID)
	}

	if err := validSignInCookie(resp); err != nil {
		t.Fatal(err)
	}

}

func Test_SignIn_Fail(t *testing.T) {
	user := User{
		ID:       1,
		Username: "cooldancer_123",
		Password: "123",
	}
	userStore := TestUserStore{user}

	sess := gumtest.NewTestSession(
		"444",
		user.Username,
		time.Now().Add(1*time.Hour),
	)
	sessionStore := TestSessionStore{sess}

	r := gin.New()
	h := SignIn(sessionStore, userStore)
	r.GET("/:name/:password", h)

	name := base64.StdEncoding.EncodeToString([]byte(user.Username))
	pass := base64.StdEncoding.EncodeToString([]byte("wrong"))

	url := fmt.Sprintf("/%v/%v", name, pass)
	resp := gumtest.NewRouter(r).ServeHTTP("GET", url, "")

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("Expect %v was %v", http.StatusUnauthorized, resp.Code)
	}

}

func validSignInCookie(r *httptest.ResponseRecorder) error {
	v, ok := r.HeaderMap["Set-Cookie"]
	if !ok {
		m := fmt.Sprintf("Expect a cookie was %v", r.HeaderMap)
		return errors.New(m)
	}
	if !strings.Contains(v[0], XSRFCookieName) {
		m := fmt.Sprintf("Expect %v was %v",
			XSRFCookieName, r.HeaderMap)
		return errors.New(m)
	}

	return nil
}

func newTestUser() User {
	return User{
		Username:  "devilXX",
		Password:  "123",
		FirstName: "Dare",
		LastName:  "Devil",
		Email:     "devil@hell.de",
		IsActive:  true,
		LastLogin: time.Now(),
	}
}

func fillTestDB(t *testing.T, db *gorp.DbMap) User {
	user := newTestUser()

	user.Password = NewSha512Password(user.Password)

	err := db.Insert(&user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}
