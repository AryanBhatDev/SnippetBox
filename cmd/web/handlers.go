package main

import (
	"errors"
	"fmt"
	_ "html/template"
	"net/http"
	"strconv"

	"github.com/AryanBhatDev/SnippetBox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// files := []string{
	//     "./ui/html/base.tmpl",
	//     "./ui/html/partials/nav.tmpl",
	//     "./ui/html/pages/home.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	//     app.serverError(w, r, err)
	//     return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	//     app.serverError(w, r, err)
	// }
}

func (app application) homePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from post"))
}

func (app application) snippetView(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	idStringToInt, err := strconv.Atoi(id)

	if idStringToInt < 0 || err != nil {

		http.NotFound(w, r)
		return
	}

	data, err := app.snippets.Get(idStringToInt)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", data)
}

func (app application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("hello", "wello")
	w.Write([]byte("Write your snippet here"))
}

func (app application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
