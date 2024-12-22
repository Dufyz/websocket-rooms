package controller

import (
	"encoding/json"
	"fmt"
	ws "socket-server/internal/infra/websocket"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type websocketController struct {
	room_manager *ws.RoomManager
}

func NewWebsocketController() *websocketController {
	return &websocketController{
		room_manager: ws.NewRoomManager(),
	}
}

func (wc *websocketController) HandleConnection(ctx echo.Context) error {
	room_manager := wc.room_manager

	wsHandler := websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()
		defer room_manager.RemoveClientFromAllRooms(conn)

		var current_room *ws.Room

		for {
			var msg string
			err := websocket.Message.Receive(conn, &msg)
			if err != nil {
				break
			}

			var event ws.Event
			if err := json.Unmarshal([]byte(msg), &event); err != nil {
				continue
			}

			fmt.Println("Received event: ", event)
			fmt.Println("Received msg: ", msg)

			switch event.Type {
			case "join":
				if current_room != nil {
					current_room.RemoveClient(conn)
				}
				current_room = wc.room_manager.GetOrCreateRoom(event.Room_id)
				current_room.AddClient(conn)

			case "leave":
				if current_room != nil {
					current_room.RemoveClient(conn)
					current_room = nil
				}

			case "message":
				if current_room != nil {
					var payload ws.MessagePayload
					payloadBytes, _ := json.Marshal(event.Payload)
					if err := json.Unmarshal(payloadBytes, &payload); err != nil {
						continue
					}

					messageEvent := ws.Event{
						Type:    "message",
						Room_id: event.Room_id,
						Payload: map[string]interface{}{
							"message": payload.Message,
						},
					}

					current_room.BroadcastEvent(messageEvent)
				}
			}

		}

	})

	wsHandler.ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}
