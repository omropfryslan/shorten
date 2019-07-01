package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var lastID int64

func decodeHandler(response http.ResponseWriter, request *http.Request, db Database) {
	short := mux.Vars(request)["short"]

	url, err := db.Get(short)
	if err != nil {
		http.Error(response, `{"error": "No such URL"}`, http.StatusNotFound)
		return
	}

	http.Redirect(response, request, url, 301)
}

func main() {
	if os.Getenv("BASE_URL") == "" {
		log.Fatal("BASE_URL environment variable must be set")
	}
	if os.Getenv("DB_PATH") == "" {
		log.Fatal("DB_PATH environment variable must be set")
	}
	if os.Getenv("API_KEY") == "" {
		log.Fatal("API_KEY environment variable must be set")
	}

	db := sqlite{Path: path.Join(os.Getenv("DB_PATH"), "db.sqlite")}
	db.Init()

	baseURL := os.Getenv("BASE_URL")
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	apiKey := os.Getenv("API_KEY")

	r := mux.NewRouter()
	r.HandleFunc("/api/save",
		func(response http.ResponseWriter, request *http.Request) {
			encodeHandler(response, request, db, baseURL)
		}).Methods("POST")

	r.HandleFunc("/api/shorten",
		func(response http.ResponseWriter, request *http.Request) {
			encodeAPIHandler(response, request, db, baseURL, apiKey)
		})

	r.HandleFunc("/{short}",
		func(response http.ResponseWriter, request *http.Request) {
			decodeHandler(response, request, db)
		})

	port := os.Getenv("PORT")
	if port == "" {
		port = "1337"
	}

	port = ":" + port

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	log.Println("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(port, handlers.LoggingHandler(os.Stdout, r)))
}
