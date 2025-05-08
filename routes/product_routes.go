package routes

import (
	"tesodev_interview/controllers"
	"tesodev_interview/middleware"

	"github.com/labstack/echo/v4"
)

func ProductRoute(e *echo.Echo) {
	e.POST("/product", middleware.LogMiddleware(controllers.CreateProduct))
}
