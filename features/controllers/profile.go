package controllers

import (
	"media-server/features/models"
	"media-server/features/views"
	"media-server/utils"

	"github.com/labstack/echo/v4"
)

func ProfilePage(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return utils.RenderTempl(c, 200, views.LoginPage("user data corrupted"))
	}

	if c.Request().Method == "POST" {
		password := c.FormValue("password")
		newpassword := c.FormValue("newpassword")

		err := models.ChangePassword(username, password, newpassword)
		if err != nil {
			return utils.RenderTempl(c, 500, views.ProfilePage(username, err.Error()))
		}

		return utils.RenderTempl(c, 200, views.ProfilePage(username))
	}

	return utils.RenderTempl(c, 200, views.ProfilePage(username))
}
