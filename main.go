package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/achristie/udemy-modern-go/config"
	"github.com/achristie/udemy-modern-go/handlers"
	"github.com/achristie/udemy-modern-go/render"
)

func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println("Listening on port 8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
