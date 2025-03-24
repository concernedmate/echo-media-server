package controllers

import (
	"media-server/features/views"
	"media-server/utils"

	"github.com/labstack/echo/v4"
)

func DashboardPage(c echo.Context) error {
	return utils.RenderTempl(c, 200, views.DashboardPage())
}
