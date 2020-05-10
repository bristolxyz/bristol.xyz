package routes

import (
	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/bristolxyz/bristol.xyz/utils"
	"github.com/labstack/echo/v4"
)

func init() {
	clients.EchoInstance.GET("/resend", func(c echo.Context) error {
		// Send the verification email.
		err := utils.SendVerificationEmail(c.Get("user").(*models.User))
		if err != nil {
			c.Response().Status = 400
			c.Response().Write([]byte(err.Error()))
			return nil
		}

		// Send 204.
		c.Response().Status = 204
		return nil
	})
}
