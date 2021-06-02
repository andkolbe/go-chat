package main

// A package is a way to group functions, and it's made up of all the files in the same directory
// A main function executes by default when you run the main package

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-chat/internal/config"
	"github.com/andkolbe/go-chat/internal/driver"
	"github.com/andkolbe/go-chat/internal/handlers"

	_ "github.com/jackc/pgx/v4/stdlib"
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

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go-chatroom user=postgres password=andrew00")
	if err != nil {
		log.Fatal("Cannot connect to db. Dying...")
	}
	log.Println("Connected to database!")

	defer db.SQL.Close() // refactor this if I switch to the run function!

	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe(":8080", mux)

	// our app config (where we can put whatever we want) and our db (a pointer to a db driver) are available to all of our handlers
	// right now our db only holds postgres, but if we change or add more in the future, that can easily be refactored
	repo := handlers.NewRepo(&app, db)
	// pass the repo variable back to the handlers
	handlers.NewHandlers(repo)
	// // gives the render component of our app access to the app config variable
	// helpers.NewHelpers(&app)

}

// func getAllRows(conn *sql.DB) error {
// 	rows, err := conn.Query("SELECT id, username, first_name, last_name, email FROM users ORDER BY id")
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	// must close connection to rows
// 	defer rows.Close()

// 	// pull something out of the rows variable and put it in GO
// 	var username, firstName, lastName, email string
// 	var id int

// 	for rows.Next() { // Next ranges over all of the rows
// 		err := rows.Scan(&id, &username, &firstName, &lastName, &email) // Scan the Query in the order they are written and store them in variables
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 		fmt.Println("Record is", id, username, firstName, lastName, email)
// 	}

// 	if err = rows.Err(); err != nil {
// 		log.Fatal("Error scanning rows", err)
// 	}

// 	fmt.Println("---------------------------------------------------")

// 	return nil
// }

