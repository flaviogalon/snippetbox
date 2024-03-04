package main

import (
	"net/http"

	"github.com/justinas/alice"
	"snippetbox.flaviogalon.github.io/internal/utils"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(
		utils.NeuteredFileSystem{Fs: http.Dir(app.appConfig.staticAssertsDir)},
	)

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	standardMiddleware := alice.New(
		app.recoverPanic,
		app.logRequests,
		secureHeaders,
	)

	return standardMiddleware.Then(mux)
}
