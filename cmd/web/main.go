package main

// A package is a way to group functions, and it's made up of all the files in the same directory
// A main function executes by default when you run the main package

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/go-chat/internal/config"
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

	// connect to db
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=go-chatroom user=postgres password=andrew00")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect: %v\n", err))
	}
	defer conn.Close() // must close connection

	log.Println("Connected to database!")

	// test my connection
	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot ping dataabse")
	}

	log.Println("Ping database")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// insert a row
	// query := `
	// 	INSERT INTO users (username, first_name, last_name, password, email)
	// 	VALUES ($1, $2, $3, $4, $5)
	// `

	// log.Println("Inserted a row!")

	// _, err = conn.Exec(query, "test", "Test", "Guy", "password", "test@gmail.com")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // get rows from table again
	// err = getAllRows(conn)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// update a row
	// stmt := `
	// 	UPDATE users  
	// 	SET username = $1
	// 	WHERE id = $2
	// `

	// log.Println("Updated a row!")

	// // Exec executes a query without returning any rows
	// _, err = conn.Exec(stmt, "bobby12", 1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// get rows from table again
	// err = getAllRows(conn)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// get one row by id
	query := `
		SELECT id, username, first_name, last_name, email 
		FROM users 
		WHERE id = $1 
	`
	var username, firstName, lastName, email string
	var id int
	row := conn.QueryRow(query, 1)
	err = row.Scan(&id, &username, &firstName, &lastName, &email)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("QueryRow returns", id, username, firstName, lastName, email)

	// delete a row
	// query = `
	// 	DELETE FROM users
	// 	WHERE id = $1
	// `
	// _, err = conn.Exec(query, 4)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Deleted a row!")

	// get rows again
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe(":8080", mux)
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("SELECT id, username, first_name, last_name, email FROM users ORDER BY id")
	if err != nil {
		log.Println(err)
		return err
	}
	// must close connection to rows
	defer rows.Close()

	// pull something out of the rows variable and put it in GO
	var username, firstName, lastName, email string
	var id int

	for rows.Next() { // Next ranges over all of the rows
		err := rows.Scan(&id, &username, &firstName, &lastName, &email) // Scan the Query in the order they are written and store them in variables
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("Record is", id, username, firstName, lastName, email)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning rows", err)
	}

	fmt.Println("---------------------------------------------------")

	return nil
}

