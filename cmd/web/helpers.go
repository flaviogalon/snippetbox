package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Write an error message and stack trace to custom error logger then sends
// a 500 response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// Send a specific status code and corresponding description to caller
func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(
	w http.ResponseWriter,
	statusCode int,
	pageName string,
	data *templateData,
) {
	ts, ok := app.templateCache[pageName]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", pageName)
		app.serverError(w, err)
		return
	}

	// Using a buffer to prevent runtime template errors from leaking to the user
	buffer := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buffer, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(statusCode)

	buffer.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
