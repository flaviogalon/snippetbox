package main

import (
	"net/http"

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

	// Middleware chain: logRequests -> secureHeaders -> ServeMux
	return app.logRequests(secureHeaders(mux))
}
