package handlers

import (
	"log"
	"net/http"

	"github.com/andkolbe/go-chat/internal/helpers"
	"github.com/andkolbe/go-chat/internal/models"
)


func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// renews the token on every new log in
	// prevents session fixation attack
	_ = m.App.Session.RenewToken(r.Context()) 

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		ClientError(w, r, http.StatusBadRequest)
	}

	id, hash, err := Repo.DB.Authenticate(r.Form.Get("username"), r.Form.Get("password"))
	if err == models.ErrInvalidCredentials {
		app.Session.Put(r.Context(), "error", "Invalid login")
		err := helpers.RenderPage(w, r, "login.jet", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
		return
	} else if err != nil {
		log.Println(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	// we authenticated. Get the user
	u, err := m.DB.GetUserByID(id)
	if err != nil {
		log.Println(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	app.Session.Put(r.Context(), "user_id", id)
	app.Session.Put(r.Context(), "hashedPassword", hash)
	app.Session.Put(r.Context(), "flash", "You've been logged in successfully")
	app.Session.Put(r.Context(), "user", u)

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}