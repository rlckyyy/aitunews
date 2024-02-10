package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/news", dynamicMiddleware.ThenFunc(app.showNews))
	mux.Post("/news/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNews))
	mux.Get("/news/creation", dynamicMiddleware.ThenFunc(app.showCreateNewsForm))
	mux.Del("/news/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteNews))
	mux.Get("/contacts", dynamicMiddleware.ThenFunc(app.showContacts))
	mux.Get("/news/students", dynamicMiddleware.ThenFunc(app.showCategoryNews))
	mux.Get("/news/staff", dynamicMiddleware.ThenFunc(app.showCategoryNews))
	mux.Get("/news/applicants", dynamicMiddleware.ThenFunc(app.showCategoryNews))

	mux.Get("/admin", dynamicMiddleware.ThenFunc(app.getAdminPage))
	mux.Get("/adminShow", dynamicMiddleware.ThenFunc(app.showUser))
	mux.Post("/updateRole", dynamicMiddleware.ThenFunc(app.updateUserRoleHandler))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
