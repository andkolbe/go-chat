package main

import (
	"net/http"

	"github.com/andkolbe/go-chat/internal/handlers"
	"github.com/go-chi/chi/v5"
)


func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/login", http.HandlerFunc(handlers.Login))

	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))


	return mux
}