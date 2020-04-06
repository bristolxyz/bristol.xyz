package main

import (
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/labstack/echo/v4"
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
		if token != "" {
			c.Set("user", models.GetUserByToken(token))
		}
		return next(c)
	}
}
