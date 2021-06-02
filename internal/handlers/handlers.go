package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/andkolbe/go-chat/internal/config"
	"github.com/andkolbe/go-chat/internal/driver"
	"github.com/andkolbe/go-chat/internal/repository"
	"github.com/andkolbe/go-chat/internal/repository/dbrepo"
)

// repository pattern allows us to swap components out of our app with minimal changes required to the code base

// the repository used by the handlers
var Repo *Repository

// the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
	}
}
// when we call newRepo, we pass it the app config (a pointer to config.AppConfig), 
// and the database connection pool (a pointer to driver.DB, which holds to db connection pool)
// we then populate the Repository type with all of the information we receive as parameters
// and hand that back as a pointer to Repository

// sets the repository for the handlers on the main package
func NewHandlers(r *Repository) {
	Repo = r
}

// every web handler in Go must have a response writer and a pointer to a request

// must have this to use the jet templating engine
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"), // look at the files in the html folder
	jet.InDevelopmentMode(), // we don't have to restart our app every time we make a change to a jet template
)

// giving the handlers a receiver links them together with the repository, so all of the handlers have access to the repository
// those handlers have access to everything inside of the app config and the database driver

// Home Page
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.html", nil)
	if err != nil {
		log.Println(err)
	}
}

// Register Page
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "register.html", nil)
	if err != nil {
		log.Println(err)
	}
}

// Chat Room Page
func (m *Repository) Chat(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "chat.html", nil)
	if err != nil {
		log.Println(err)
	}
}

// renders pages
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}