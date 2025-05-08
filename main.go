package main

import (
	"fmt"
	"os"
	"tesodev_interview/configs"
	"tesodev_interview/middleware"
	"tesodev_interview/routes"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func init() {
	if os.Getenv("APP_ENV") == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.LogMiddleware)

	configs.ConnectDB()
	fmt.Println("db bağlanma sonrası")

	routes.ProductRoute(e)

	e.Logger.Fatal(e.Start(":3000"))
}
