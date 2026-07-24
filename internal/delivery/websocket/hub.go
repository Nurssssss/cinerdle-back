package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu    sync.Mutex
	queue []*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{}
}
func (hu *Hub) AddToQueue(conn *websocket.Conn) {
	hu.mu.Lock()
	defer hu.mu.Unlock()
	hu.queue = append(hu.queue, conn)

}

func (hu *Hub) TryMatch() {
	hu.mu.Lock()
	defer hu.mu.Unlock()
	if len(hu.queue) >= 2 {
		firstPlayer := hu.queue[0]
		secondPlayer := hu.queue[1]
		firstPlayer.WriteMessage(websocket.TextMessage, []byte("game_found"))
		secondPlayer.WriteMessage(websocket.TextMessage, []byte("game_found"))
		hu.queue = hu.queue[2:]
	}

}
