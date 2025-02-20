package main

import (
	"fmt"
	"media-server/configs"
	"media-server/features/routes"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "[INFO ] ${time_custom} - method=${method}, uri=${uri}, status=${status}\n",
		CustomTimeFormat: time.DateTime,
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			currentTime := time.Now().Local().Format(time.DateTime)
			fmt.Printf("[ERROR] %s - %s\n", currentTime, err.Error())
			if configs.LOG_STACK {
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

	routes.AuthRoutes(e)
	routes.DashboardRoutes(e)

	e.Logger.Fatal(e.Start(":3000"))
}
