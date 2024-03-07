package main

import (
	"errors"
	"fmt"
	"net/http"
	"relucky.net/aitunews/pkg/models"
	"strconv"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	newsList, err := app.news.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		NewsList: newsList,
	})
}

func (app *application) showNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	n, err := app.news.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	comments, err := app.news.GetComments(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	userID := app.session.GetInt(r, "authenticatedUserID")
	app.render(w, r, "show.page.tmpl", &templateData{
		News:      n,
		Comments:  comments,
		SessionId: userID,
	})
}

func (app *application) createNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
	}
	if user.Role != models.RoleTeacher && user.Role != models.RoleAdmin {
		app.clientError(w, http.StatusForbidden)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.FormValue("category")

	if title == "" || content == "" || category == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if len(title) > 20 || len(content) > 200 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	validCategories := []string{"Students", "Staff", "Applicants"}

	if !contains(validCategories, category) {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// go to logic layer
	id, err := app.news.Insert(title, content, category)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", id), http.StatusSeeOther)
}

func (app *application) showCreateNewsForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{})
}

func (app *application) showCategoryNews(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Path[len("/news/"):]

	if category == "" {
		http.NotFound(w, r)
		return
	}
	category = strings.Title(category)

	newsList, err := app.news.GetByCategory(category)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "category.page.tmpl", &templateData{
		Category: category,
		NewsList: newsList,
	})
}

func (app *application) showContacts(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "contacts.page.tmpl", &templateData{})
}

func (app *application) deleteNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.news.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func contains(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}

// user auth part

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Println(name + email + password)
	role := models.RoleUser
	err = app.users.Insert(name, email, password, role)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			app.session.Put(r, "flash", "Email already exists.")
			app.render(w, r, "signup.page.tmpl", nil)
			return
		}
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	id, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.session.Put(r, "flash", "Invalid email or password.")
			app.render(w, r, "login.page.tmpl", nil)
			return
		}
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) getAdminPage(w http.ResponseWriter, r *http.Request) {
	userList, err := app.users.GetUsers()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "admin.page.tmpl", &templateData{
		UserList: userList,
	})
}

func (app *application) showUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "showUser.page.tmpl", &templateData{
		User: user,
	})
}

func (app *application) updateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID := r.PostForm.Get("userID")

	newRole := r.PostForm.Get("newRole")

	err = app.users.UpdateUserRole(userID, newRole)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/adminShow?id=%s", userID), http.StatusSeeOther)
}

// comments section

func (app *application) addComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	newsID, err := strconv.Atoi(r.FormValue("newsID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID := app.session.GetInt(r, "authenticatedUserID")
	text := r.FormValue("text")
	err = app.comments.Insert(text, userID, newsID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", newsID), http.StatusSeeOther)
}

func (app *application) deleteComment(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)

	commentID, err := strconv.Atoi(r.FormValue("commentID"))
	if err != nil || commentID < 1 {
		app.serverError(w, err)
		return
	}
	newsId, err := app.comments.GetNewsId(commentID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	authorId, err := app.comments.GetAuthorId(commentID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if user.Role != "admin" && userID != authorId {
		app.session.Put(r, "flash", "You can only delete your own comments!")
		http.Redirect(w, r, fmt.Sprintf("/news?id=%d", newsId), http.StatusSeeOther)
		return
	}
	err = app.comments.Delete(commentID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", newsId), http.StatusSeeOther)
}
