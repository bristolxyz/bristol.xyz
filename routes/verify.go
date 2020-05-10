package routes

import (
	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/labstack/echo/v4"
)

func init() {
	clients.EchoInstance.GET("/verify/:id", func(c echo.Context) error {
		// Defer redirect.
		defer c.Redirect(302, "/")

		// Get the URL param.
		ID := c.Param("id")

		// Verify the user if possible.
		v := models.GetVerificationCode(ID)
		if v != nil {
			v.Verify()
		}

		// Return nil.
		return nil
	})
}
