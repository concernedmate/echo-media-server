package routes

import (
	"github.com/labstack/echo/v4"
)

func ResourcesRoutes(e *echo.Echo) {
	e.Static("/assets", "./resources/assets")
}
