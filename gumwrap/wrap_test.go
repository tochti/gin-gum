package gumwrap

import (
	"database/sql"
	"testing"

	"gopkg.in/gorp.v1"

	"github.com/gin-gonic/gin"
)

func TestSQLDB(t *testing.T) {
	db := &sql.DB{}
	SQLDB(func(c *gin.Context, db *sql.DB) {
	}, db)
}

func TestGorp(t *testing.T) {
	db := &gorp.DbMap{}
	Gorp(func(c *gin.Context, db *gorp.DbMap) {
	}, db)
}
