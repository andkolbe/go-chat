package main

// A package is a way to group functions, and it's made up of all the files in the same directory
// a package is a collection of common source code files
// A main function executes by default when you run the main package

// two types of packages. Executable and reusable
// executable - generates a file that we can run. Only the main package
// reusable - code used as 'helpers'. Good place to put reusable logic

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-chat/internal/config"
	"github.com/andkolbe/go-chat/internal/driver"
	"github.com/andkolbe/go-chat/internal/handlers"
	"github.com/andkolbe/go-chat/internal/models"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close() // db won't close until the main function stops running

	mux := routes()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe(":8080", mux)

}

func run() (*driver.DB, error) {

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
	go handlers.ListenToWsChannel()

	// our app config (where we can put whatever we want) and our db (a pointer to a db driver) are available to all of our handlers
	// right now our db only holds postgres, but if we change or add more in the future, that can easily be refactored
	repo := handlers.NewRepo(&app, db)
	// pass the repo variable back to the handlers
	handlers.NewHandlers(repo)
	// // gives the render component of our app access to the app config variable
	// helpers.NewHelpers(&app)

	return db, nil
}



