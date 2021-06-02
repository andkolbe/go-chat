package main

import (
	"net/http"

	"github.com/andkolbe/go-chat/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	// middleware allows you process a request as it comes into your web app and perform some action on it
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/register", http.HandlerFunc(handlers.Register))
	mux.Get("/chat", http.HandlerFunc(handlers.Chat))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))

	// if a user is disconnected, and then reconnects, they rejoin automatically
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}