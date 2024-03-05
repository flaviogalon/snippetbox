package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.flaviogalon.github.io/internal/models"
)

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippetModel.Latest()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippets = snippets

	app.render(
		w,
		http.StatusOK,
		"home.tmpl.html",
		templateData,
	)
}

// Display a single snippet handler
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippet = snippet

	app.render(
		w,
		http.StatusOK,
		"view.tmpl.html",
		templateData,
	)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

// Create a snippet handler
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	formData := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	// Data validation
	if strings.TrimSpace(title) == "" {
		formData.FieldErrors["title"] = "This field can't be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		formData.FieldErrors["title"] = "This field can't be more than 100 characters long"
	}

	if strings.TrimSpace(content) == "" {
		formData.FieldErrors["content"] = "This field can't be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		formData.FieldErrors["expires"] = "This field must be equal 1, 7 or 365"
	}

	// If data validation failed render the page with the errors
	if len(formData.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippetModel.Insert(
		formData.Title,
		formData.Content,
		formData.Expires,
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
