package middlewares

import (
	"media-server/features/models"
	"media-server/features/views"
	"media-server/utils"

	"github.com/labstack/echo/v4"
)

func CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Cookie("token")

		token, err := c.Cookie("token")
		if err != nil {
			return utils.RenderTempl(c, 200, views.LoginPage("token invalid"))
		}

		username, err := models.CheckToken(token.Value)
		if err != nil {
			return utils.RenderTempl(c, 200, views.LoginPage(err.Error()))
		}

		c.Set("username", username)
		return next(c)
	}
}
