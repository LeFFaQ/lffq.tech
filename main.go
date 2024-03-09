package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", AboutHandler)

	mux.HandleFunc("GET /blog/", BlogHandler)
	mux.HandleFunc("GET /blog/{id}/", BlogByIdHandler)

	http.Handle("/static", http.FileServer(http.Dir("./static/")))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("./template/about.gohtml")

	t.Execute(w, nil)

}

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/blog.gohtml")

	t.Execute(w, nil)
}

func BlogByIdHandler(w http.ResponseWriter, r *http.Request) {

}
