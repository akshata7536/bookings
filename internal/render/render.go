package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/akshata7536/bookings/internal/config"
	"github.com/akshata7536/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template

	var err error
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Println("error while creating template cache", err)
		}
	}
	for key := range tc {
		log.Println("Available template in cache:", key)
	}
	log.Printf("%v", tc)
	t, ok := tc[tmpl]
	if !ok {
		log.Fatalf("Template '%s' not found in template cache. Available templates: %v", tmpl, tc)
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// log.Println("name ::", name)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
