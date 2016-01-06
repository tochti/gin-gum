package gumrest

import "github.com/gin-gonic/gin"

type ErrorMessage struct {
	Message string `json:"message"`
}

// Creates following JSON Response
//
// {
//   "message": "..."
// }
func ErrorResponse(c *gin.Context, code int, err error) {
	c.JSON(code, ErrorMessage{
		Message: err.Error(),
	})
}
