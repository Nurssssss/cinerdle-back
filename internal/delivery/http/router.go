package http

import (
	"cinerdle-back/internal/delivery/websocket"
	"net/http"
)

type Router struct {
	authHandler          *AuthHandler
	authMiddleware       *AuthMiddleware
	movieHandler         *MovieHandler
	gameHandler          *GameHandler
	gameWebscoketHandler *websocket.GameWebsocketHandler
}

func NewRouter(authHandler *AuthHandler, authMiddleware *AuthMiddleware, movieHandler *MovieHandler, gameHandler *GameHandler, gameWebscoketHandler *websocket.GameWebsocketHandler) *Router {
	return &Router{
		authHandler:          authHandler,
		authMiddleware:       authMiddleware,
		movieHandler:         movieHandler,
		gameHandler:          gameHandler,
		gameWebscoketHandler: gameWebscoketHandler,
	}
}

func (rt *Router) Setup(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", rt.authHandler.Register)
	mux.HandleFunc("POST /login", rt.authHandler.Login)
	mux.HandleFunc("POST /verify-email", rt.authHandler.VerifyEmail)
	mux.Handle("GET /get-profile", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.authHandler.GetProfile)))
	mux.Handle("GET /search-movie", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.movieHandler.SearchMovie)))
	mux.Handle("POST /find-game", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.gameHandler.CreateGame)))
	mux.Handle("GET /handle-connection", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.gameWebscoketHandler.HandleConnection)))
}
