package main

import (
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UserMiddleware is used to handle the authentication of users where this is required.
func UserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token string
		tc, _ := c.Cookie("token")
		if tc == nil {
			token = c.Request().Header.Get("Token-Auth")
		} else {
			token = tc.Value
		}
		if token == "" {
			var b *models.User
			c.Set("user", b)
		} else {
			u := models.GetUserByToken(token)
			if u != nil {
				c.Set("token", token)
				if u.Banned {
					c.SetCookie(&http.Cookie{
						Name:   "token",
						MaxAge: -1,
					})
					u = nil
				}
			}
			c.Set("user", u)
		}
		return next(c)
	}
}
