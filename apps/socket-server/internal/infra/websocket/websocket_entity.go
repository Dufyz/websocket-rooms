package websocket

import (
	"sync"

	"golang.org/x/net/websocket"
)

type MessagePayload struct {
	Message string `json:"message"`
}

type Event struct {
	Type    string      `json:"type"`
	Room_id string      `json:"room_id"`
	Payload interface{} `json:"payload"`
}

type Room struct {
	id        string
	clients   map[*websocket.Conn]bool
	broadcast chan Event
	mutex     sync.RWMutex
}

type RoomPayload struct {
	Room_id string `json:"room_id"`
}

type RoomManager struct {
	rooms map[string]*Room
	mutex sync.RWMutex
}
