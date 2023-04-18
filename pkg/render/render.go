package render

import (
	"bytes"
	"github.com/jhowilbur/go-web-app-reservations/pkg/config"
	"github.com/jhowilbur/go-web-app-reservations/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

var applicationConfig *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	applicationConfig = a
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	// get the template cache from the app config
	var templateCache map[string]*template.Template

	if applicationConfig.UseCache {
		// template caching
		templateCache = applicationConfig.TemplateCache
	} else {
		// build template cache manually
		templateCache, _ = CreateTemplate()
	}

	// get requested template from cache
	templateWanted, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buffer := new(bytes.Buffer) // execute values get from map and execute that directly
	templateData = AddDefaultData(templateData)
	_ = templateWanted.Execute(buffer, templateData) // to help identify if it has an error and show the template

	// render the template
	_, err := buffer.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser")
	}

}

func CreateTemplate() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.TemplateData)
	myCache := map[string]*template.Template{} // same functionality like above.

	// get all the files named *.page.tmpl from ./tempaltes
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)                             // to get the last element in the path (file tmpl)
		templateSet, err := template.New(name).ParseFiles(page) // populate name with HTML

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}
