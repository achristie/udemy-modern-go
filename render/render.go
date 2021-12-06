package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/achristie/udemy-modern-go/config"
	"github.com/achristie/udemy-modern-go/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			return err
		}
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	str := fmt.Sprintf("%s/*.page.tmpl", pathToTemplates)
	pages, err := filepath.Glob(str)
	if err != nil {
		return cache, err
	}

	for _, pg := range pages {
		name := filepath.Base(pg)
		ts, err := template.New(name).Funcs(functions).ParseFiles(pg)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return cache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return cache, err
			}

			cache[name] = ts
		}
	}
	return cache, nil
}
