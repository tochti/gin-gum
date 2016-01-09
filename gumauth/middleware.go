package gumauth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumrest"
	"github.com/tochti/session-store"
)

const (
	SessionKey       = "Session"
	XSRFCookieName   = "XSRF-TOKEN"
	TokenHeaderField = "X-XSRF-TOKEN"
)

var (
	CookieErr          = errors.New("Cookie error")
	SessionNotFoundErr = errors.New("Session not found")
	HeaderNotFoundErr  = errors.New("Header not found")
	CookieNotFoundErr  = errors.New("Cookie not found")
)

func ReadSession(c *gin.Context) (s2tore.Session, error) {
	v, ok := c.Get(SessionKey)
	if !ok {
		return nil, CookieErr
	}

	s, ok := v.(s2tore.Session)
	if !ok {
		return nil, CookieErr
	}

	return s, nil
}

func SignedIn(s s2tore.SessionStore) func(gin.HandlerFunc) gin.HandlerFunc {
	return func(h gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			ok := Bouncer(c, s)
			if !ok {
				return
			}

			h(c)
		}

	}
}

func Bouncer(c *gin.Context, s s2tore.SessionStore) bool {
	token := c.Request.Header.Get(TokenHeaderField)
	if token == "" {
		cookie, err := c.Request.Cookie(XSRFCookieName)
		if err != nil {
			gumrest.ErrorResponse(
				c,
				http.StatusUnauthorized,
				CookieNotFoundErr,
			)
			return false
		}
		token = cookie.Value
		if token == "" {
			gumrest.ErrorResponse(
				c,
				http.StatusUnauthorized,
				HeaderNotFoundErr,
			)
		}
	}

	session, ok := s.ReadSession(token)
	if !ok {
		gumrest.ErrorResponse(
			c,
			http.StatusUnauthorized,
			SessionNotFoundErr,
		)
		return false

	}

	c.Set(SessionKey, session)
	return true
}
