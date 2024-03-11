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

	// Unprotected application routes
	dynamicMid := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamicMid.ThenFunc(app.home))
	// Snippet
	router.Handler(
		http.MethodGet,
		"/snippet/view/:id",
		dynamicMid.ThenFunc(app.snippetView),
	)
	router.Handler(
		http.MethodGet,
		"/user/signup",
		dynamicMid.ThenFunc(app.userSignup),
	)
	router.Handler(
		http.MethodPost,
		"/user/signup",
		dynamicMid.ThenFunc(app.userSignupPost),
	)
	router.Handler(
		http.MethodGet,
		"/user/login",
		dynamicMid.ThenFunc(app.userLogin),
	)
	router.Handler(
		http.MethodPost,
		"/user/login",
		dynamicMid.ThenFunc(app.userLoginPost),
	)

	// Protected application routes
	protected := dynamicMid.Append(app.requireAuthentication)
	router.Handler(
		http.MethodGet,
		"/snippet/create",
		protected.ThenFunc(app.snippetCreate),
	)
	router.Handler(
		http.MethodPost,
		"/snippet/create",
		protected.ThenFunc(app.snippetCreatePost),
	)
	router.Handler(
		http.MethodPost,
		"/user/logout",
		protected.ThenFunc(app.userLogoutPost),
	)

	standardMiddleware := alice.New(
		app.recoverPanic,
		app.logRequests,
		secureHeaders,
	)

	return standardMiddleware.Then(router)
}
