package main

import (
	"database/sql"
	"github.com/gragas/woobloo-frontend/config"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"net/http"
)

var db *sql.DB
var templates map[string]*template.Template

func main() {
	// connect to the database
	var err error
	config := config.LoadConfig("config.json")
	db, err = sql.Open(config.DBType, config.ConnectionString())
	if err != nil {
		panic(err)
	}

	// initialize maps, etc.
	templates = make(map[string]*template.Template)

	// setup static handlers
	const (
		cssPath    = "/static/css/"
		imagesPath = "/static/images/"
	)
	http.Handle(cssPath,
		http.StripPrefix(cssPath, http.FileServer(http.Dir("."+cssPath))))
	http.Handle(imagesPath,
		http.StripPrefix(imagesPath, http.FileServer(http.Dir("."+imagesPath))))

	// setup dynamic handlers
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(config.ServerAddr, nil)
}

// Handlers
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}
	data := struct {
		Title      string
		TestString string
	}{Title: "Woobloo Civilizations", TestString: "Hello, World!"}
	getTemplate("html/index.html").Execute(w, data)
}

func getTemplate(htmlPath string) *template.Template {
	var t *template.Template
	if templates[htmlPath] == nil {
		// if the template is not in the cache, generate it
		bytes, err := ioutil.ReadFile(htmlPath)
		if err != nil {
			panic(err)
		}
		t, err = template.New(htmlPath).Parse(string(bytes))
		if err != nil {
			panic(err)
		}
		// and add it to the cache
		templates[htmlPath] = t
	} else {
		// otherwise, grab it from the cache
		t = templates[htmlPath]
	}
	return t
}
