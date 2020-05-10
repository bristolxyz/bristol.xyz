package main

import (
	"github.com/bristolxyz/bristol.xyz/models"
	"os"

	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/bristolxyz/bristol.xyz/routes"
)

func main() {
	// Load .env if we can.
	_ = godotenv.Load()

	// Load Sentry.
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		panic(err)
	}

	// Load MongoDB.
	err = clients.CreateMongoClient()
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	// Run post-database load tasks.
	models.RunDatabaseInitTasks()

	// Load S3.
	clients.S3Init()

	// Enable the default logging and recovery middleware.
	clients.EchoInstance.Use(middleware.Logger())
	clients.EchoInstance.Use(middleware.Recover())

	// Create the Sentry handler.
	clients.EchoInstance.Use(sentryecho.New(sentryecho.Options{}))

	// Add the user middleware.
	clients.EchoInstance.Use(UserMiddleware)

	// Serve the static folder.
	clients.EchoInstance.Static("/static", "static")

	// Start the web server.
	clients.EchoInstance.Logger.Fatal(clients.EchoInstance.Start(":8080"))
}
