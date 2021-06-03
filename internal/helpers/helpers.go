package helpers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/andkolbe/go-chat/internal/models"
	"github.com/justinas/nosurf"
)

// must have this to use the jet templating engine
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"), // look at the files in the html folder
	jet.InDevelopmentMode(), // we don't have to restart our app every time we make a change to a jet template
)

// DefaultData adds default data which is accessible to all templates
func DefaultData(td models.TemplateData, r *http.Request, w http.ResponseWriter) models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	// td.IsAuthenticated = IsAuthenticated(r)
	// td.PreferenceMap = app.PreferenceMap
	// // if logged in, store user id in template data
	// if td.IsAuthenticated {
	// 	u := app.Session.Get(r.Context(), "user").(models.User)
	// 	td.User = u
	// }

	// td.Flash = app.Session.PopString(r.Context(), "flash")
	// td.Warning = app.Session.PopString(r.Context(), "warning")
	// td.Error = app.Session.PopString(r.Context(), "error")

	return td
}

// renders pages
func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	// add default template data
	var td models.TemplateData
	if data != nil {
		td = data.(models.TemplateData)
	}

	// add default data
	td = DefaultData(td, r, w)

	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, vars, td)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}