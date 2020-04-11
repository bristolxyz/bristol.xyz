package routes

import (
	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/labstack/echo/v4"
)

func init() {
	clients.EchoInstance.GET("/", func(c echo.Context) error {
		_, err := c.Response().Write([]byte("Hello World!"))
		return err
	})
}
