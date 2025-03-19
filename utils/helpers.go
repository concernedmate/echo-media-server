package utils

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type ResponseFormat struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func RenderTempl(c echo.Context, httpStatus int, view templ.Component) error {
	c.Response().Writer.WriteHeader(httpStatus)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return view.Render(c.Request().Context(), c.Response().Writer)
}

func ResponseJSON(c echo.Context, httpStatus int, message string, data any) error {
	return c.JSON(httpStatus, ResponseFormat{
		Message: message,
		Data:    data,
	})
}
