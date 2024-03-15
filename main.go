package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/v60/github"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var GH_TOKEN = os.Getenv("TOKEN")
var md = goldmark.New(
	goldmark.WithExtensions(
		&frontmatter.Extender{},
	),
)

var CACHED_POSTS []BlogContent

func init() {
	pctx := parser.NewContext()

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: GH_TOKEN})
	tc := oauth2.NewClient(ctx, ts)

	gh := github.NewClient(tc)

	_, directory, _, err := gh.Repositories.GetContents(ctx, "LeFFaQ", "blog-content", "", nil)

	if err != nil {
		log.Print("Something wrong with initial request to repository")
		fmt.Println(err)
	}

	bc := BlogContent{}

	for _, content := range directory {
		if content.GetType() == "dir" {
			slug := content.GetPath()

			file, _, _, _ := gh.Repositories.GetContents(ctx, "LeFFaQ", "blog-content", slug+"/readme.md", nil)
			fileContent, _ := file.GetContent()

			meta := Meta{}

			buf := new(bytes.Buffer)
			if err := md.Convert([]byte(fileContent), buf, parser.WithContext(pctx)); err != nil {
				log.Print("Something wrong with parsing md to html")
				log.Fatal(err)
			}

			fm := frontmatter.Get(pctx)
			if fm == nil {
				println("no frontmatter found")
			}
			if err := fm.Decode(&meta); err != nil {
				log.Print("Something wrong with parsing frontmatter")
				log.Fatal(err)
			}
			bc.Slug = slug
			bc.Meta = meta
			bc.Content = buf.String()

			CACHED_POSTS = append(CACHED_POSTS, bc)
		}
	}

	// Test purpose
	for _, post := range CACHED_POSTS {
		fmt.Printf("Post %s \n", post.Slug)
		fmt.Printf("Title: %s \n", post.Meta.Title)
		fmt.Printf("Content: %s \n", post.Content)
	}
}

func main() {
	fmt.Println("App is running!")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", AboutHandler)

	mux.HandleFunc("GET /blog/", BlogHandler)
	mux.HandleFunc("GET /blog/{id}/", BlogByIdHandler)

	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	log.Fatal(http.ListenAndServe(":8080", mux))

}

func AboutHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles(
		"./template/index.gohtml",
		"./template/about.gohtml",
		"./template/header.gohtml") // I don't know how to DRY it

	err := t.ExecuteTemplate(w, "index", nil)

	if err != nil {
		fmt.Println(err)
	}
}

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	//data := BlogData{}

	mapped := map[string]Meta{}
	for _, post := range CACHED_POSTS {
		mapped[post.Slug] = post.Meta
		fmt.Printf("%s, %s \n", post.Slug, post.Meta.Title)
	}

	t, _ := template.ParseFiles(
		"./template/index.gohtml",
		"./template/blog.gohtml",
		"./template/header.gohtml") // I don't know how to DRY it

	_ = t.ExecuteTemplate(w, "index", mapped)

}

func BlogByIdHandler(w http.ResponseWriter, r *http.Request) {

}

type Meta struct {
	Title string    `yaml:"title"`
	Desc  string    `yaml:"description"`
	Cover string    `yaml:"cover"`
	Date  time.Time `yaml:"date"`
	Tags  []string  `yaml:"tags"`
}

// BlogContent Describes parsed content of blog, it's metadata and slug. Used for "cache"
type BlogContent struct {
	Slug    string
	Meta    Meta
	Content string
}
