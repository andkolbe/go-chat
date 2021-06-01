package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

// every web handler in go must have a response writer and a pointer to a request

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(), // we don't have to restart our app every time we make a change to a jet template
)

// Home Page
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.html", nil)
	if err != nil {
		log.Println(err)
	}
}

// Login Page
func Register(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "register.html", nil)
	if err != nil {
		log.Println(err)
	}
}

// Chat Room Page
func Chat(w http.ResponseWriter, r *http.Request) {
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