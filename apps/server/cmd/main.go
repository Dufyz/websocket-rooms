package main

import (
	"fmt"
	"socket-server/internal/interfaces/rest/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	routes.UseRoutes(e)

	if err := e.Start(":3000"); err != nil {
		fmt.Println("Error starting server")
	}
}
