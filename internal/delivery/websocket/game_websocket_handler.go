package websocket

import (
	"cinerdle-back/internal/domain"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type GameWebsocketHandler struct {
	gameUseCase GameUseCase
	hub         *Hub
}

type GameUseCase interface {
	CreateGame(firstPlayer, secondPlayer uuid.UUID) (*domain.GameSession, error)
}

func NewWebsocketHandler(gameUseCase GameUseCase, hub *Hub) *GameWebsocketHandler {
	return &GameWebsocketHandler{
		gameUseCase: gameUseCase,
		hub:         hub,
	}
}

func (weh *GameWebsocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	weh.hub.AddToQueue(conn)
	weh.hub.TryMatch()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}

}
