package routes

import (
	"media-server/features/controllers"

	"github.com/labstack/echo/v4"
)

func DashboardRoutes(e *echo.Echo) {
	e.GET("/dashboard", controllers.DashboardPage)
}
