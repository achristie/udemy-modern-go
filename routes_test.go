package main

import (
	"testing"

	"github.com/achristie/udemy-modern-go/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// passed
	default:
		t.Errorf("type is no *chi.Mux, type is %T", v)
	}
}
