package routes

import (
	"media-server/features/controllers"

	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Echo) {
	e.GET("/", controllers.LoginPage)
	e.POST("/", controllers.LoginPage)
}
