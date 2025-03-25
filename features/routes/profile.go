package routes

import (
	"media-server/features/controllers"
	"media-server/features/middlewares"

	"github.com/labstack/echo/v4"
)

func ProfileRoutes(e *echo.Echo) {
	e.GET("/profile", controllers.ProfilePage, middlewares.CheckAuth)
	e.POST("/profile", controllers.ProfilePage, middlewares.CheckAuth)
}
