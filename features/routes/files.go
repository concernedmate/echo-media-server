package routes

import (
	"media-server/features/controllers/api"

	"github.com/labstack/echo/v4"
)

func APIFilesRoutes(e *echo.Echo) {
	e.GET("/api/v1/files", api.ListFile)
	e.GET("/api/v1/files/download", api.DownloadFile)
	e.POST("/api/v1/files/upload", api.UploadFile)
}
