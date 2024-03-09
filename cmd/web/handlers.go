package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.flaviogalon.github.io/internal/models"
	"snippetbox.flaviogalon.github.io/internal/validator"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
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

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	formData := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// Data validation
	formData.CheckField(validator.NotBlank(formData.Title), "title", "This field can't be blank")
	formData.CheckField(
		validator.MaxChars(formData.Title, 100),
		"title",
		"This field can't be more than 100 characters long",
	)
	formData.CheckField(validator.NotBlank(formData.Content), "content", "This field can't be blank")
	formData.CheckField(
		validator.PermittedInt(formData.Expires, 1, 7, 365),
		"expires",
		"This field must be equal 1, 7 or 365",
	)

	if !formData.Valid() {
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
