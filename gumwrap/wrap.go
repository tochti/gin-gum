package gumwrap

import (
	"database/sql"

	"gopkg.in/gorp.v1"

	"github.com/gin-gonic/gin"
)

// Pass a sql db connection to a func(*gin.Context, *sql.DB)
func SQLDB(h func(*gin.Context, *sql.DB), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c, db)
	}
}

func Gorp(h func(*gin.Context, *gorp.DbMap), db *gorp.DbMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c, db)
	}
}
