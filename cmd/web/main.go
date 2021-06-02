package main

// A package is a way to group functions, and it's made up of all the files in the same directory
// A main function executes by default when you run the main package

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-chat/internal/config"
	"github.com/andkolbe/go-chat/internal/handlers"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

	// enable sessions in the main package
	session = scs.New()
	session.Lifetime = 24 * time.Hour // active for 24 hours
	// stores the session in cookies by default. Can switch to Redis
	session.Cookie.Persist = true // cookie persists when the browser window is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // makes sure the cookies are encrypted and use https. Change to true for production

	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe(":8080", mux)
}

