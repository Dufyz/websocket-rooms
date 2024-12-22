package routes

import (
	"socket-server/internal/interfaces/rest/controller"

	"github.com/labstack/echo/v4"
)

func UseRoutes(e *echo.Echo) {
	api := e.Group("/api")

	WebsocketController := controller.NewWebsocketController()
	api.GET("/web-socket", WebsocketController.HandleConnection)
}
