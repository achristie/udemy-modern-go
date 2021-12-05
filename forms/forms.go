package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Form struct {
	url.Values
	Errors errors
}

var validate *validator.Validate

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks
func (f *Form) Has(field string, r *http.Request) bool {
	x := f.Get(field)

	if x == "" {
		f.Errors.Add(field, "this field cannot be blank")
		return false
	}

	return true
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	validate = validator.New()
	errs := validate.Var(field, "email")
	if errs != nil {
		f.Errors.Add(field, fmt.Sprint(errs))
		return
	}
}
