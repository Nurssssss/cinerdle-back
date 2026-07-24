package websocket

import (
	"math/rand"
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
		index := rand.Intn(2)
		randomWhoStarted := hu.queue[index]
		randomWhoWaited := hu.queue[1-index]
		randomWhoStarted.WriteMessage(websocket.TextMessage, []byte("you_start_first"))
		randomWhoWaited.WriteMessage(websocket.TextMessage, []byte("you_wait"))

		hu.queue = hu.queue[2:]
	}

}
