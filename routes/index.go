package routes

import (
	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/bristolxyz/bristol.xyz/utils"
	"github.com/labstack/echo/v4"
)

func init() {
	clients.EchoInstance.GET("/", func(c echo.Context) error {
		// Get the user.
		User := c.Get("user").(*models.User)

		// Return the rendered content.
		c.Response().Header().Set("Content-Type", "text/html;charset=utf-8")
		_, err := c.Response().Write(utils.GenerateBase("Home", "The homepage for Bristol.xyz - Bristol's leading and open news site",
			"", "<h1>Hello World!</h1>", User))
		return err
	})
}
