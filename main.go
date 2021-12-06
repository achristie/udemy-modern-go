package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/achristie/udemy-modern-go/config"
	"github.com/achristie/udemy-modern-go/handlers"
	"github.com/achristie/udemy-modern-go/helpers"
	"github.com/achristie/udemy-modern-go/models"
	"github.com/achristie/udemy-modern-go/render"
	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var session *scs.SessionManager
var errorLog *log.Logger
var infoLog *log.Logger

func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on port 8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {

	gob.Register(models.Reservation{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalf("Cannot create template cache: %s", err)
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)
	return nil
}
