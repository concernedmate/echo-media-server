package controllers

import (
	"fmt"
	"media-server/configs"
	"media-server/features/views"
	"media-server/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginPage(c echo.Context) error {
	if c.Request().Method == "POST" {
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s/dashboard", configs.BASE_URL))
	}

	return utils.RenderTempl(c, 200, views.LoginPage())
}
