package api

import (
	"fmt"
	"media-server/configs"
	"media-server/features/models"
	"media-server/features/websockets"
	"media-server/utils"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

func CreateFolder(c echo.Context) error {
	dir := c.FormValue("dir")
	if dir == "" || dir == "/" || strings.Contains(dir, "/") {
		return utils.ResponseJSON(c, 400, "[Bad Request] invalid directory name", nil)
	}

	curr_dir := c.FormValue("curr_dir")
	if curr_dir == "" || curr_dir == "/" {
		dir = "/" + dir
	} else {
		dir = curr_dir + "/" + dir
	}

	username, ok := c.Get("username").(string)
	if !ok {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] invalid username", nil)
	}

	err := models.CreateFolder(dir, username)
	if err != nil {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] "+err.Error(), nil)
	}

	return utils.ResponseJSON(c, 200, "[Success]", nil)
}

func DeleteFolder(c echo.Context) error {
	var req struct {
		Directory string `json:"directory" validate:"required"`
	}

	err := c.Bind(&req)
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}

	if req.Directory == "" {
		return utils.ResponseJSON(c, 400, "[Bad Request] directory is required", nil)
	}

	err = models.DeleteFolder(req.Directory)
	if err != nil {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] "+err.Error(), nil)
	}

	return utils.ResponseJSON(c, 200, "[Success]", nil)
}

func UploadMultipleFiles(c echo.Context) error {
	multipart, err := c.MultipartForm()
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}
	dir := multipart.Value["dir"][0]
	if dir == "" {
		dir = "/"
	}

	username, ok := c.Get("username").(string)
	if !ok {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] invalid username", nil)
	}

	err = models.UploadMultipleFiles(multipart.File["files"], dir, username)
	if err != nil {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] "+err.Error(), nil)
	}

	return utils.ResponseJSON(c, 200, "[Success]", nil)
}

func UploadMultipleFileViaWS(c echo.Context) error {
	err := websockets.UploadMultipleFiles(c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
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

	file, err := os.Open(path.Join(configs.UPLOAD_BASEDIR(), metadata.FileId))
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}
	defer file.Close()

	return c.Attachment(path.Join(configs.UPLOAD_BASEDIR(), metadata.FileId), metadata.Filename)
}

func ShowContent(c echo.Context) error {
	file_id := c.QueryParam("file_id")
	if file_id == "" {
		return utils.ResponseJSON(c, 400, "[Bad Request] file_id is required", nil)
	}

	metadata, err := models.GetFileMetadata(file_id)
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}

	split := strings.Split(metadata.Filename, ".")
	ext := split[len(split)-1]

	switch ext {
	case "png", "jpg", "gif":
		file, err := os.Open(path.Join(configs.UPLOAD_BASEDIR(), metadata.FileId))
		if err != nil {
			return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
		}
		defer file.Close()

		return c.Inline(path.Join(configs.UPLOAD_BASEDIR(), metadata.FileId), metadata.Filename)
	case "mp4":
		// TODO
	}

	// send no preview image
	file, err := os.Open("./resources/assets/static/images/no_preview_available.png")
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}
	defer file.Close()

	return c.Inline("./resources/assets/static/images/no_preview_available.png", file.Name())
}

func DeleteFile(c echo.Context) error {
	var req struct {
		FileID string `json:"file_id" validate:"required"`
	}

	err := c.Bind(&req)
	if err != nil {
		return utils.ResponseJSON(c, 400, "[Bad Request] "+err.Error(), nil)
	}

	if req.FileID == "" {
		return utils.ResponseJSON(c, 400, "[Bad Request] file_id is required", nil)
	}

	err = models.DeleteFile(req.FileID)
	if err != nil {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] "+err.Error(), nil)
	}

	return utils.ResponseJSON(c, 200, "[Success]", nil)
}

func ListFile(c echo.Context) error {
	dir := c.FormValue("dir")
	if dir == "" {
		dir = "/"
	}

	username, ok := c.Get("username").(string)
	if !ok {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] "+`invalid user data`, nil)
	}

	files, err := models.ListFiles(username, dir)
	if err != nil {
		return utils.ResponseJSON(c, 500, "[Internal Server Error] "+err.Error(), nil)
	}

	return utils.ResponseJSON(c, 200, "[Success]", files)
}
