package main

import (
	"fmt"
	"log"
	"media-server/configs"
	"media-server/features/models"
	"media-server/features/routes"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		log.Fatal(err.Error())
	}
	if err := models.InitSQLite(); err != nil {
		log.Fatal(err.Error())
	}

	e := echo.New()

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieDomain:   "",
		CookieSecure:   false,
		CookieHTTPOnly: false,
		CookieSameSite: http.SameSiteStrictMode,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "[INFO ] ${time_custom} - method=${method}, uri=${uri}, status=${status}\n",
		CustomTimeFormat: time.DateTime,
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			currentTime := time.Now().Local().Format(time.DateTime)
			fmt.Printf("[ERROR] %s - %s\n", currentTime, err.Error())
			if configs.LOG_STACK() {
				fmt.Printf("[STACK] ")
				for _, data := range stack {
					if string(data) == "\t" {
						continue
					}
					if string(data) == "\n" {
						fmt.Printf("\n[STACK] ")
					} else {
						fmt.Printf("%s", string(data))
					}
				}
			}
			return c.JSON(500, echo.Map{"message": "Internal Server Error"})
		},
	}))
	e.Use()

	routes.ResourcesRoutes(e)

	routes.AuthRoutes(e)
	routes.DashboardRoutes(e)
	routes.DriveRoutes(e)
	routes.ProfileRoutes(e)

	routes.APIFilesRoutes(e)

	routes.ErrorRoutes(e)

	e.Logger.Fatal(e.Start(":3000"))
}
