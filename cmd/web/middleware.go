package main

import (
	"net/http"
	"github.com/justinas/nosurf"
)

// adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// uses cookies to make sure the token it generates for us is on a per page basis
	csrfHandler.SetBaseCookie(http.Cookie {
		HttpOnly: true,
		Path: "/", // "/" means apply this cookie to the entire site
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// web servers are not state aware, so we need to add middleware that tells this application to remember state using sessions
// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}