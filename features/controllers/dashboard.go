package controllers

import (
	"media-server/features/views"
	"media-server/utils"

	"github.com/labstack/echo/v4"
)

func DashboardPage(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return utils.RenderTempl(c, 200, views.LoginPage("user data corrupted"))
	}
	return utils.RenderTempl(c, 200, views.DashboardPage(username))
}
