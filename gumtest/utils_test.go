package gumtest

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestTestSession(t *testing.T) {
	exp := time.Now().Add(1 * time.Hour)
	s := TestSession{
		token:   "123",
		userID:  "12",
		expires: exp,
	}

	if s.Token() != "123" {
		t.Fatalf("Expect 123 was %v", s.Token())
	}

	if s.UserID() != "12" {
		t.Fatalf("Expect 12 was %v", s.UserID())
	}

	if s.Expires() != exp {
		t.Fatalf("Expect %v was %v", exp, s.Expires())
	}
}

func TestTestHandler(t *testing.T) {
	handlerDone := false

	th := func(c *gin.Context) {
		if id := c.Param("id"); id != "1" {
			t.Fatalf("Expect 1 was %v", id)
		}

		body := make([]byte, 4)
		_, err := c.Request.Body.Read(body)
		if err != nil {
			t.Fatal(err)
		}

		if string(body) != "body" {
			t.Fatalf("Expect body was %v", body)
		}

		handlerDone = true
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/:id", th)

	NewRouter(r).ServeHTTP("GET", "/1", "body")

	if !handlerDone {
		t.Fatal("Expect to successfully run test handler")
	}
}

func TestMockAuther(t *testing.T) {
	handlerDone := false

	th := func(c *gin.Context) {
		s := c.MustGet("Session").(TestSession)
		if id := s.UserID(); id != "23" {
			t.Fatalf("Expect 23 was %v", id)
		}

		if id := c.Param("id"); id != "1" {
			t.Fatalf("Expect 1 was %v", id)
		}

		handlerDone = true
	}

	userID := "23"
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/:id", MockAuther(th, userID))

	NewRouter(r).ServeHTTP("GET", "/1", "")

	if !handlerDone {
		t.Fatal("Expect to successfully run test handler")
	}
}

func TestJSONResponse(t *testing.T) {
	data := struct {
		Data string
	}{
		Data: "test",
	}
	resp := JSONResponse{200, data}

	expectBody := "{\n\t\"Data\": \"test\"\n}"

	if resp.String() != expectBody {
		t.Fatalf("Expect %v was %v", expectBody, resp.String())
	}

	if resp.Code() != 200 {
		t.Fatal("Expect 200 was", resp.Code())
	}
}

func TestEqualJSONResponse(t *testing.T) {
	writer := httptest.NewRecorder()
	writer.HeaderMap.Add("Content-Type", "application/json; charset=utf-8")

	data := struct {
		Data string
	}{
		Data: "test",
	}

	respBody, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	writer.Write(respBody)

	expect := JSONResponse{200, data}

	if err := EqualJSONResponse(expect, writer); err != nil {
		t.Fatal("Expect nil was", err)
	}
}

func TestEqualJSONResponse_Fail(t *testing.T) {
	writer := httptest.NewRecorder()

	data := struct {
		Data string
	}{
		Data: "test",
	}
	expect := JSONResponse{200, data}

	if err := EqualJSONResponse(expect, writer); err == nil {
		t.Fatal("Expect error was nil")
	}
}

func TestIfaceSlice(t *testing.T) {
	result := IfaceSlice([]string{"a", "b"})

	expect := []interface{}{
		"a",
		"b",
	}

	for i, s := range expect {
		if s != result[i] {
			t.Fatalf("Expect %v was %v", s, result[i])
		}
	}
}
