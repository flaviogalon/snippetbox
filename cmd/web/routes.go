package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"snippetbox.flaviogalon.github.io/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Embedded File Server
	fileServer := http.FileServer(http.FS(ui.Files))

	router.Handler(
		http.MethodGet,
		"/static/*filepath",
		fileServer,
	)
	router.HandlerFunc(
		http.MethodGet,
		"/ping",
		ping,
	)

	// Unprotected application routes
	dynamicMid := alice.New(
		app.sessionManager.LoadAndSave,
		noSurf,
		app.authenticate,
	)

	router.Handler(http.MethodGet, "/", dynamicMid.ThenFunc(app.home))
	router.Handler(
		http.MethodGet,
		"/about",
		dynamicMid.ThenFunc(app.about),
	)
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

	// Protected application routes (copying all mid from unprotected)
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
	router.Handler(
		http.MethodGet,
		"/account/view",
		protected.ThenFunc(app.accountView),
	)

	standardMiddleware := alice.New(
		app.recoverPanic,
		app.logRequests,
		secureHeaders,
	)

	return standardMiddleware.Then(router)
}
