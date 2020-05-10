package routes

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/bristolxyz/bristol.xyz/utils"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

func init() {
	// Get the register template.
	b, err := ioutil.ReadFile("./templates/register.html")
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	t, err := template.New("register").Parse(string(b))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	// Handles the register GET request.
	clients.EchoInstance.GET("/register", func(c echo.Context) error {
		// Get the user.
		User := c.Get("user").(*models.User)

		// If the user is not nil, redirect home and return.
		if User != nil {
			return c.Redirect(307, "/")
		}

		// Create the rendered register page.
		b := bytes.Buffer{}
		err := t.Execute(&b, nil)
		if err != nil {
			return err
		}

		// Return the rendered content.
		c.Response().Header().Set("Content-Type", "text/html;charset=utf-8")
		_, err = c.Response().Write(utils.GenerateBase("Register", "The registration page for Bristol.xyz.",
			"", b.String(), User))
		return err
	})

	// Handles the register POST request.
	clients.EchoInstance.POST("/register", func(c echo.Context) error {
		// The render function for the register page.
		RenderLoginPageWithError := func(Message string) error {
			// Create the rendered register page.
			b := bytes.Buffer{}
			err := t.Execute(&b, Message)
			if err != nil {
				return err
			}

			// Return the rendered content.
			c.Response().Header().Set("Content-Type", "text/html;charset=utf-8")
			_, err = c.Response().Write(utils.GenerateBase("Register", "The registration page for Bristol.xyz.",
				"", b.String(), c.Get("user").(*models.User)))
			return err
		}

		// Get all elements from the form.
		Email := c.FormValue("email")
		Password := c.FormValue("password")
		ConfirmPassword := c.FormValue("confirmPassword")
		FirstName := c.FormValue("firstName")
		LastName := c.FormValue("lastName")
		FirstNamePtr := &FirstName
		LastNamePtr := &LastName

		// If any required options are blank, return a error.
		if Email == "" {
			return RenderLoginPageWithError("E-mail address field is blank.")
		}
		if FirstName == "" {
			FirstNamePtr = nil
		}
		if LastName == "" {
			LastNamePtr = nil
		}
		if Password == "" {
			return RenderLoginPageWithError("Password field is blank.")
		}
		if ConfirmPassword == "" {
			return RenderLoginPageWithError("Confirm password field is blank.")
		}

		// Check if the confirmation matches the password.
		if ConfirmPassword != Password {
			return RenderLoginPageWithError("Confirmation password does not match.")
		}

		// Make the e-mail address lower case.
		Email = strings.ToLower(Email)

		// Check if a user already exists with this e-mail address.
		if models.GetUserByEmail(Email) != nil {
			return RenderLoginPageWithError("A user already exists with this e-mail address.")
		}

		// Create a new user.
		u, token := models.NewUser(Email, Password, false, false, FirstNamePtr, LastNamePtr)
		c.SetCookie(&http.Cookie{
			Name:   "token",
			Value:  token,
			MaxAge: 2592000,
		})

		// Send the confirmation e-mail.
		utils.SendVerificationEmail(u)

		// Redirect home.
		return c.Redirect(302, "/")
	})
}
