package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", AboutHandler)

	mux.HandleFunc("GET /blog/", BlogHandler)
	mux.HandleFunc("GET /blog/{id}/", BlogByIdHandler)

	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	log.Fatal(http.ListenAndServe(":8080", mux))

}

func AboutHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("./template/index.gohtml", "./template/about.gohtml")

	err := t.ExecuteTemplate(w, "index", nil)

	if err != nil {
		fmt.Println(err)
	}
}

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/blog.gohtml", "./template/index.html")

	log.Fatal(t.Execute(w, nil))
}

func BlogByIdHandler(w http.ResponseWriter, r *http.Request) {

}
