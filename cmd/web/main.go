package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/jhowilbur/go-web-app-reservations/pkg/config"
	"github.com/jhowilbur/go-web-app-reservations/pkg/handlers"
	"github.com/jhowilbur/go-web-app-reservations/pkg/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

// calling template cache from config
const portNUmber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false
	app.UseCache = true

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := render.CreateTemplate()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = templateCache

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// handle endpoints
	http.Handle("/metrics", promhttp.Handler()) // Prometheus metrics
	handlers.RecordMetrics()                    // Prometheus custom metrics

	//http.HandleFunc("/", repo.Home)
	//http.HandleFunc("/about", repo.About)

	log.Println(fmt.Sprintf("Server starting on port %s", portNUmber))
	//_ = http.ListenAndServe(portNUmber, nil)

	serve := &http.Server{
		Addr:    portNUmber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}
