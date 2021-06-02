package main

import (
	"net/http"

	"github.com/andkolbe/go-chat/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/register", http.HandlerFunc(handlers.Register))
	mux.Get("/chat", http.HandlerFunc(handlers.Chat))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))

	// if a user is disconnected, and then reconnects, they rejoin automatically
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}