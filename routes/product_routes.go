package routes

import (
	"tesodev_interview/handlers"
	"tesodev_interview/middleware"

	"github.com/labstack/echo/v4"
)

func ProductRoute(e *echo.Echo) {
	e.POST("/product", middleware.LogMiddleware(handlers.CreateProduct))
	e.PUT("/product/:product_id", middleware.LogMiddleware(handlers.UpdateProduct))
	e.GET("/product/:product_id", middleware.LogMiddleware(handlers.GetAProduct))
	e.GET("/product", middleware.LogMiddleware(handlers.GetAllProduct))
	e.DELETE("/product/:product_id", middleware.LogMiddleware(handlers.DeleteProduct))
	e.PATCH("/product/:product_id", middleware.LogMiddleware(handlers.UpdateSingleFeild))
}

func SearchRoute(e *echo.Echo) {
	e.GET("/search", middleware.LogMiddleware(handlers.SearchProducts))
}
