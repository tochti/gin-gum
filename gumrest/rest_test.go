package gumrest

import (
	"encoding/json"
	"errors"
	"strings"
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

	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/", "")

	expect, err := json.Marshal(ErrorMessage{
		Message: "fuck off",
	})
	if err != nil {
		t.Fatal(err)
	}

	result := strings.Replace(resp.Body.String(), "\n", "", -1)
	if string(expect) != result {
		t.Fatalf("Expect (%v) was (%v)", string(expect), result)
	}
}
