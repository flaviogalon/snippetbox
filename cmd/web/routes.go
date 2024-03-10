package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir(app.appConfig.staticAssertsDir))
	router.Handler(
		http.MethodGet,
		"/static/*filepath",
		http.StripPrefix("/static", fileServer),
	)

	dynamicMid := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamicMid.ThenFunc(app.home))
	router.Handler(
		http.MethodGet,
		"/snippet/view/:id",
		dynamicMid.ThenFunc(app.snippetView),
	)
	router.Handler(
		http.MethodGet,
		"/snippet/create",
		dynamicMid.ThenFunc(app.snippetCreate),
	)
	router.Handler(
		http.MethodPost,
		"/snippet/create",
		dynamicMid.ThenFunc(app.snippetCreatePost),
	)

	standardMiddleware := alice.New(
		app.recoverPanic,
		app.logRequests,
		secureHeaders,
	)

	return standardMiddleware.Then(router)
}
