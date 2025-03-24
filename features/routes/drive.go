package routes

import (
	"media-server/features/controllers"
	"media-server/features/middlewares"

	"github.com/labstack/echo/v4"
)

func DriveRoutes(e *echo.Echo) {
	e.GET("/drive", controllers.DrivePage, middlewares.CheckAuth)
}
