package api

import (
	"media-server/features/models"
	"media-server/utils"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

func UploadFile(c echo.Context) error {
	dir := c.FormValue("dir")
	if dir == "" {
		dir = "/"
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}

	err = models.UploadFile(file, dir, "admin")
	if err != nil {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] "+err.Error(), nil)
	}

	return utils.ResponseJSON(c, 200, "[Success]", nil)
}

func DownloadFile(c echo.Context) error {
	file_id := c.QueryParam("file_id")
	if file_id == "" {
		return utils.ResponseJSON(c, 400, "[Bad Request] file_id is required", nil)
	}

	metadata, err := models.GetFileMetadata(file_id)
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}

	file, err := os.Open(path.Join("./uploads", metadata.FileId))
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}
	defer file.Close()

	return c.Attachment(path.Join("./uploads", metadata.FileId), metadata.Filename)
}

func ListFile(c echo.Context) error {
	dir := c.FormValue("dir")
	if dir == "" {
		dir = "/"
	}

	files, err := models.ListFiles(dir)
	if err != nil {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] "+err.Error(), nil)
	}

	return utils.ResponseJSON(c, 200, "[Success]", files)
}
