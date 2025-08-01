package routes

import (
	"media-server/features/controllers/api"
	"media-server/features/middlewares"

	"github.com/labstack/echo/v4"
)

func APIFilesRoutes(e *echo.Echo) {
	e.GET("/api/v1/files", api.ListFile, middlewares.CheckAuth)
	e.GET("/api/v1/files/download", api.DownloadFile, middlewares.CheckAuth)
	e.GET("/api/v1/files/show", api.ShowContent, middlewares.CheckAuth)
	e.POST("/api/v1/files/upload/batch", api.UploadMultipleFiles, middlewares.CheckAuth)
	e.DELETE("/api/v1/files/delete", api.DeleteFile, middlewares.CheckAuth)

	e.POST("/api/v1/files/folder/create", api.CreateFolder, middlewares.CheckAuth)
	e.DELETE("/api/v1/files/folder/delete", api.DeleteFolder, middlewares.CheckAuth)

	e.GET("/ws/v1/files/upload/batch", api.UploadMultipleFileViaWS, middlewares.CheckAuth)
}
