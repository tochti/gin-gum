package gumauth

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumrest"
	"github.com/tochti/session-store"
)

const (
	NameRequestField = "name"
	PassRequestField = "password"
)

var (
	SignInErr     = errors.New("Wrong password")
	UserExistsErr = errors.New("User already exists")
)

type (
	Session struct {
		Token   string    `json:"token"`
		UserID  string    `json:"user_id"`
		Expires time.Time `json:"expires"`
	}
)

func SQLNewUser(db *sql.DB, u *User) error {
	q := fmt.Sprintf("SELECT id FROM %v WHERE username=?", UserTable)
	tmp := -1
	err := db.QueryRow(q, u.Username).Scan(&tmp)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if tmp != -1 {
		return UserExistsErr
	}

	q = fmt.Sprintf(`
		INSERT 
		INTO %v
		(username,
		 first_name,
		 last_name,
		 email,
		 password,
		 is_active,
		 last_login)
		VALUES(?,?,?,?,?,?,?)
		`, UserTable)

	res, err := db.Exec(q,
		u.Username,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		u.IsActive,
		u.LastLogin,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func CreateUserSQL(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := User{}
		err := c.BindJSON(&user)
		if err != nil {
			gumrest.ErrorResponse(
				c,
				http.StatusNotAcceptable,
				err,
			)
			return
		}

		user.IsActive = true
		user.LastLogin = time.Now()
		user.Password = NewSha512Password(user.Password)

		err = SQLNewUser(db, &user)
		if err != nil {
			gumrest.ErrorResponse(
				c,
				http.StatusNotAcceptable,
				err,
			)
			return
		}

		c.JSON(http.StatusCreated, user)
	}

}

func SignIn(s s2tore.SessionStore, u UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		Signer(c, s, u)
	}
}

func Signer(c *gin.Context, s s2tore.SessionStore, u UserStore) {
	nameParam := c.Params.ByName(NameRequestField)
	tmp, err := base64.StdEncoding.DecodeString(nameParam)
	if err != nil {
		gumrest.ErrorResponse(c, http.StatusNotAcceptable, err)
		return
	}
	name := string(tmp)
	passParam := c.Params.ByName(PassRequestField)
	tmp, err = base64.StdEncoding.DecodeString(passParam)
	if err != nil {
		gumrest.ErrorResponse(c, http.StatusNotAcceptable, err)
		return
	}
	pass := string(tmp)

	user, err := u.FindUser(name)
	if err != nil {
		gumrest.ErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	if !user.ValidPassword(pass) {
		gumrest.ErrorResponse(c, http.StatusUnauthorized, SignInErr)
		return
	}

	expire := time.Now().Add(24 * time.Hour)
	session, err := s.NewSession(user.ID(), expire)
	if err != nil {
		gumrest.ErrorResponse(c, http.StatusNotAcceptable, err)
		return
	}

	c.Set(SessionKey, session)

	cookie := http.Cookie{
		Name:    XSRFCookieName,
		Value:   session.Token(),
		Expires: session.Expires(),
		// Setze Path auf / ansonsten kann angularjs
		// diese Cookie nicht finden und in späteren
		// Request nicht mitsenden.
		Path: "/",
	}
	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusAccepted, Session{
		Token:   session.Token(),
		UserID:  session.UserID(),
		Expires: session.Expires(),
	})
}
