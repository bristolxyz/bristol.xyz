package routes

import (
	"context"
	"net/http"

	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	clients.EchoInstance.GET("/logout", func(c echo.Context) error {
		// If user is nil, just return the site redirect.
		if c.Get("user").(*models.User) == nil {
			c.Redirect(302, "/")
			return nil
		}

		// Get the token.
		Token := c.Get("token").(string)

		// Delete the token.
		_, err := clients.MongoDatabase.Collection("tokens").DeleteOne(context.TODO(), bson.M{"_id": Token})
		if err != nil {
			c.Redirect(302, "/")
			return err
		}
		c.SetCookie(&http.Cookie{
			Name:   "token",
			MaxAge: -1,
		})
		c.Redirect(302, "/")
		return nil
	})
}
