package main

import (
	"os"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load Sentry.
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		panic(err)
	}

	// Create the web server.
	e := echo.New()

	// Enable the default logging and recovery middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create the handler.
	e.Use(sentryecho.New(sentryecho.Options{}))

	// Test route.
	e.GET("/", func(c echo.Context) error {
		_, err := c.Response().Write([]byte("Hello World!"))
		return err
	})

	// Start the web server.
	e.Logger.Fatal(e.Start(":8080"))
}
