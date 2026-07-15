package http

import (
	"net/http"
)

type Router struct {
	authHandler    *AuthHandler
	authMiddleware *AuthMiddleware
	movieHandler   *MovieHandler
}

func NewRouter(authHandler *AuthHandler, authMiddleware *AuthMiddleware, movieHandler *MovieHandler) *Router {
	return &Router{
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
		movieHandler:   movieHandler,
	}
}

func (rt *Router) Setup(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", rt.authHandler.Register)
	mux.HandleFunc("POST /login", rt.authHandler.Login)
	mux.HandleFunc("POST /verify-email", rt.authHandler.VerifyEmail)
	mux.Handle("GET /get-profile", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.authHandler.GetProfile)))
	mux.Handle("GET /search-movie", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.movieHandler.SearchMovie)))
}
