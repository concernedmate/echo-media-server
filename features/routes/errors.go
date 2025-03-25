package routes

import (
	"media-server/features/views/errorpages"
	"media-server/utils"

	"github.com/labstack/echo/v4"
)

func ErrorRoutes(e *echo.Echo) {
	e.GET("/*", func(c echo.Context) error {
		return utils.RenderTempl(c, 200, errorpages.NotFound404())
	})

}
