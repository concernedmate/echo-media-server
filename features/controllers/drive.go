package controllers

import (
	"fmt"
	"media-server/features/models"
	"media-server/features/views"
	"media-server/utils"

	"github.com/labstack/echo/v4"
)

func DrivePage(c echo.Context) error {
	dir := c.QueryParam("dir")
	if dir == "" {
		dir = "/"
	}

	username, ok := c.Get("username").(string)
	if !ok {
		return utils.RenderTempl(c, 200, views.DrivePage([]models.FileMetadata{}, `invalid user data`))
	}
	files, err := models.ListFiles(username, dir)
	if err != nil {
		return utils.RenderTempl(c, 200, views.DrivePage([]models.FileMetadata{}, fmt.Sprintf(`error getting files data: %s`, err.Error())))
	}

	return utils.RenderTempl(c, 200, views.DrivePage(files))
}
