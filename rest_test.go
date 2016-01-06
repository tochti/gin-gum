package gumrest

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumtest"
)

func TestErrorResponse(t *testing.T) {
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		err := errors.New("fuck off")
		ErrorResponse(c, 404, err)
	})

	resp := gumtest.TestRouter(r).ServeHTTP("GET", "/", "")

	expect, err := json.Marshal(ErrorMessage{
		Message: "fuck off",
	})

	if expect != resp.Body.String() {
		t.Fatalf("Expect %v was %v", expect, resp.Body.String())
	}
}
