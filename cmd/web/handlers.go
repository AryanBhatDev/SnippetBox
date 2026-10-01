package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app application) home(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/pages/base.tmpl",
		"./ui/html/pages/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)

	if err != nil {
		app.serverError(w, r, err)
	}
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

	fmt.Fprintf(w, "Your snippet of id: %v", id)
}

func (app application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("hello", "wello")
	w.Write([]byte("Write your snippet here"))
}

func (app application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("hello", "wello")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Post your snippet here"))
}
