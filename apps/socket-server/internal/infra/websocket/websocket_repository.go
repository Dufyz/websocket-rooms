package websocket

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

func (r *Room) handleMessages() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic in HandleMessages: ", err)
		}
	}()

	for event := range r.broadcast {
		func() {
			r.mutex.RLock()
			defer r.mutex.RUnlock()

			event_json, err := json.Marshal(event)
			if err != nil {
				fmt.Println("Error marshalling event: ", err)
				return
			}

			disconnected_clients := make([]*websocket.Conn, 0)
			for client := range r.clients {
				if err := websocket.Message.Send(client, string(event_json)); err != nil {
					fmt.Println("Error sending message: ", err)
					disconnected_clients = append(disconnected_clients, client)
				}
			}

			if len(disconnected_clients) > 0 {
				r.mutex.RUnlock()
				r.mutex.Lock()
				for _, client := range disconnected_clients {
					delete(r.clients, client)
					client.Close()
				}
				r.mutex.Unlock()
				r.mutex.RLock()
			}
		}()
	}
}

func (r *Room) AddClient(ws *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.clients[ws] = true
}

func (r *Room) RemoveClient(ws *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.clients, ws)
}

func (r *Room) BroadcastEvent(event Event) {
	select {
	case r.broadcast <- event:
		fmt.Printf("Event broadcasted to room %s\n", r.id)
	default:
		fmt.Printf("Warning: broadcast channel full for room %s\n", r.id)
	}
}

func (rm *RoomManager) GetOrCreateRoom(room_id string) *Room {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if room, ok := rm.rooms[room_id]; ok {
		return room
	}

	room := &Room{
		id:        room_id,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Event, 100),
	}
	rm.rooms[room_id] = room

	go room.handleMessages()

	return room
}

func (rm *RoomManager) RemoveClientFromAllRooms(ws *websocket.Conn) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	for _, room := range rm.rooms {
		room.RemoveClient(ws)
	}

	for room_id, room := range rm.rooms {
		if len(room.clients) == 0 {
			close(room.broadcast)
			delete(rm.rooms, room_id)
		}
	}
}
