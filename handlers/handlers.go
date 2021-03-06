package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/achristie/udemy-modern-go/config"
	"github.com/achristie/udemy-modern-go/forms"
	"github.com/achristie/udemy-modern-go/helpers"
	"github.com/achristie/udemy-modern-go/models"
	"github.com/achristie/udemy-modern-go/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyRes models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyRes

	render.RenderTemplate(w, r, "make-reservations.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.IsEmail(r.Form.Get("email"))

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservations.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}
func (m *Repository) PostSearch(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start is %s, end is %s", start, end)))
	// render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "blah",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	stringMap := make(map[string]string)
	stringMap["text"] = "andrew"
	stringMap["remoteIP"] = remoteIP
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) ReservationHandler(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
