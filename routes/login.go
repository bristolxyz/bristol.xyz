package routes

import (
	"bytes"
	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/bristolxyz/bristol.xyz/utils"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"html/template"
	"io/ioutil"
	"net/http"
)

func init() {
	// Get the login template.
	b, err := ioutil.ReadFile("./templates/login.html")
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	t, err := template.New("login").Parse(string(b))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	// Handles the login GET request.
	clients.EchoInstance.GET("/login", func(c echo.Context) error {
		// Get the user.
		User := c.Get("user").(*models.User)

		// If the user is not nil, redirect home and return.
		if User != nil {
			return c.Redirect(307, "/")
		}

		// Create the rendered login page.
		b := bytes.Buffer{}
		err := t.Execute(&b, nil)
		if err != nil {
			return err
		}

		// Return the rendered content.
		c.Response().Header().Set("Content-Type", "text/html;charset=utf-8")
		_, err = c.Response().Write([]byte(utils.GenerateBase("Login", "The login page for Bristol.xyz.",
			"", b.String(), User)))
		return err
	})

	// Handles the login POST request.
	clients.EchoInstance.POST("/login", func(c echo.Context) error {
		// The render function for the login page.
		RenderLoginPageWithError := func(Message string) error {
			// Create the rendered login page.
			b := bytes.Buffer{}
			err := t.Execute(&b, Message)
			if err != nil {
				return err
			}

			// Return the rendered content.
			c.Response().Header().Set("Content-Type", "text/html;charset=utf-8")
			_, err = c.Response().Write([]byte(utils.GenerateBase("Login", "The login page for Bristol.xyz.",
				"", b.String(), c.Get("user").(*models.User))))
			return err
		}

		// Get the e-mail and password from the form.
		Email := c.FormValue("email")
		Password := c.FormValue("password")

		// If either are blank, return a error.
		if Email == "" {
			return RenderLoginPageWithError("E-mail address field is blank.")
		}
		if Password == "" {
			return RenderLoginPageWithError("Password field is blank.")
		}

		// Handle user authentication.
		_, token := models.LoginUser(Email, Password)
		if token == nil {
			return RenderLoginPageWithError("E-mail address or password is incorrect.")
		}
		c.SetCookie(&http.Cookie{
			Name:       "token",
			Value:      *token,
			MaxAge:		2592000,
		})

		// Redirect home.
		return c.Redirect(302, "/")
	})
}
