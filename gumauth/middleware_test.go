package gumauth

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumtest"
	"github.com/tochti/session-store"
)

type (
	TestSessionStore struct {
		Session gumtest.TestSession
	}

	TestUserStore struct {
		User User
	}
)

func Test_SignedIn_OK(t *testing.T) {
	testSess := gumtest.NewTestSession(
		"123",
		"lovemaster_XXX",
		time.Now().Add(1*time.Hour),
	)
	sessStore := TestSessionStore{testSess}

	// Test if session key in gin context
	afterAuth := func(c *gin.Context) {
		sess, err := ReadSession(c)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(testSess, sess) {
			t.Fatalf("Expect %v was %v", testSess, sess)
		}

		c.JSON(http.StatusOK, nil)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	signedIn := SignedIn(sessStore)
	r.GET("/", signedIn(afterAuth))

	head := TokenHeader("123")
	resp := gumtest.NewRouter(r).ServeHTTPWithHeader("GET", "/", "", head)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expect %v was %v", http.StatusOK, resp.Code)
	}

}

func Test_SignedIn_Fail(t *testing.T) {
	sessStore := TestSessionStore{gumtest.TestSession{}}

	afterAuth := func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	}

	signedIn := SignedIn(sessStore)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/", signedIn(afterAuth))

	head := TokenHeader("1")
	resp := gumtest.NewRouter(r).ServeHTTPWithHeader("GET", "/", "", head)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("Expect %v was %v", http.StatusUnauthorized, resp.Code)
	}

}

func (s TestSessionStore) NewSession(id string, e time.Time) (s2tore.Session, error) {
	return gumtest.NewTestSession(
		"1234",
		id,
		e,
	), nil
}

func (s TestSessionStore) ReadSession(t string) (s2tore.Session, bool) {
	return s.Session, s.Session.Token() == t
}

func (s TestSessionStore) RemoveSession(t string) error {
	return nil
}

func (s TestSessionStore) RemoveExpiredSessions() (int, error) {
	return 0, nil
}

func (s TestUserStore) FindUser(n string) (StoreUser, error) {
	return NewStoreUser(s.User), nil
}

func TokenHeader(t string) http.Header {
	h := http.Header{}
	h.Add(TokenHeaderField, t)
	return h
}

func ExistsToken(tokens []string, t string) bool {
	for _, e := range tokens {
		if t == e {
			return false
		}
	}

	return true
}
