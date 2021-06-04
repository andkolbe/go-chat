package handlers

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-chat/internal/driver"
	"github.com/andkolbe/go-chat/internal/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)


var session *scs.SessionManager

func getRoutes() http.Handler {
	// what I am going to put in the session 
	gob.Register(models.User{})

	// change this to true when in production
	app.InProduction = false

	// enable sessions in the main package
	session = scs.New()
	session.Lifetime = 24 * time.Hour // active for 24 hours
	// stores the session in cookies by default. Can switch to Redis
	session.Cookie.Persist = true // cookie persists when the browser window is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // makes sure the cookies are encrypted and use https. Change to true for production

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go-chatroom user=postgres password=andrew00")
	if err != nil {
		log.Fatal("Cannot connect to db. Dying...")
	}
	log.Println("Connected to database!")

	log.Println("Starting channel listener")
	go ListenToWsChannel()

	// our app config (where we can put whatever we want) and our db (a pointer to a db driver) are available to all of our handlers
	// right now our db only holds postgres, but if we change or add more in the future, that can easily be refactored
	repo := NewRepo(app, db)
	// pass the repo variable back to the handlers
	NewHandlers(repo)
	// // gives the render component of our app access to the app config variable
	// helpers.NewHelpers(&app)

	mux := chi.NewRouter()

	// middleware allows you process a request as it comes into your web app and perform some action on it
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Login)
	mux.Post("/", Repo.PostLogin)
	mux.Get("/register", Repo.Register)
	mux.Get("/chat", Repo.Chat)
	mux.Get("/ws", WsEndPoint)

	// if a user is disconnected, and then reconnects, they rejoin automatically
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

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