package controllers

import (
	"fmt"
	"media-server/configs"
	"media-server/features/models"
	"media-server/features/views"
	"media-server/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginPage(c echo.Context) error {
	if c.Request().Method == "POST" {
		username := c.FormValue("username")
		password := c.FormValue("password")

		token, err := models.Auth(username, password, configs.JWT_SECRET())
		if err != nil {
			return utils.RenderTempl(c, 200, views.LoginPage(err.Error()))
		}

		c.SetCookie(&http.Cookie{Name: "token", Value: token})

		return c.Redirect(http.StatusFound, fmt.Sprintf("%s/dashboard", configs.BASE_URL()))
	}

	return utils.RenderTempl(c, 200, views.LoginPage())
}
