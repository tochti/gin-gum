// Some helper functions for testing gin applications

package gumtest

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const SessionKey = "Session"

type (
	TestRouter struct {
		Router *gin.Engine
	}

	// TestSession implements github.com/tochti/session-store SessionStore Interface
	TestSession struct {
		token   string
		userID  string
		expires time.Time
	}

	JSONResponse struct {
		C int
		D interface{}
	}
)

func (s TestSession) Token() string {
	return s.token
}

func (s TestSession) UserID() string {
	return s.userID
}

func (s TestSession) Expires() time.Time {
	return s.expires
}

// Perform a Request for a given gin router.
func (t *TestRouter) ServeHTTP(method, path, body string) *httptest.ResponseRecorder {
	buf := bytes.NewBufferString(body)
	req, _ := http.NewRequest(method, path, buf)

	w := httptest.NewRecorder()

	t.Router.ServeHTTP(w, req)

	return w
}

func NewRouter(r *gin.Engine) *TestRouter {
	return &TestRouter{r}
}

// Helper for testing handlers which requires authentication.
// Append the gin context with a "Session" key which
// contains a SessionStore.
func MockAuther(h gin.HandlerFunc, userID string) gin.HandlerFunc {
	t, err := NewSessionToken()

	if err != nil {
		log.Fatal(err)
	}

	sess := TestSession{
		userID:  userID,
		token:   t,
		expires: time.Now().Add(1 * time.Hour),
	}

	return func(c *gin.Context) {
		c.Set(SessionKey, sess)
		h(c)
	}

}

// Returns a Random Session Token (sha256)
func NewSessionToken() (string, error) {
	buf := make([]byte, 2)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	c := sha256.New()
	hash := fmt.Sprintf("%x", c.Sum(buf))

	return hash, nil
}

func (r JSONResponse) Code() int {
	return r.C
}

func (r JSONResponse) String() string {
	j, err := json.Marshal(r.D)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, j, "", "\t")

	return out.String()
}

// Test if content-type, http status and json body the same
func EqualJSONResponse(expect JSONResponse, rr *httptest.ResponseRecorder) error {
	if expect.Code() != rr.Code {
		m := fmt.Sprintf("Expect %v was %v", expect.Code(), rr.Code)
		return errors.New(m)
	}

	contentType := "application/json; charset=utf-8"
	if t := rr.HeaderMap.Get("Content-Type"); t != contentType {
		m := fmt.Sprintf("Expect %v was %v", contentType, t)
		return errors.New(m)
	}

	var respBody bytes.Buffer
	json.Indent(&respBody, rr.Body.Bytes(), "", "\t")

	respBodyStr := strings.Trim(respBody.String(), "\n")
	if expect.String() != respBodyStr {
		m := fmt.Sprintf("Expect %v was (%v)",
			expect.String(), respBodyStr)
		return errors.New(m)
	}

	return nil
}

// Make []"any type" to []interface{}
func IfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		log.Fatal("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// Give back now without nano seconds for better testing due tue db times
func SimpleNow() time.Time {
	tmp := time.Now()

	return time.Date(tmp.Year(), tmp.Month(), tmp.Day(),
		tmp.Hour(), tmp.Minute(), tmp.Second(), 0, tmp.Location())
}
