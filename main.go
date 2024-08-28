package main

import (
	"encoding/json"
	"fmt"
	art "gallantone.com/main/articles"
	auth "gallantone.com/main/auth"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

// handleRequests handles requests
func handleRequests() {

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("GallantOne Server is starting...") 

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", startSomewhere)
		r.Get("/webhose", serveNews)
		r.Get("/favicon.ico", favicon)
	})

	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3600,
		Debug:            false,
	})

	handler := c.Handler(r)

	// Using port 8080, use whatever port you wish - for SSL use http.ListenAndServeTLS
	log.Fatal(http.ListenAndServe(":8080", handler))   
}


// serveNews serves the news
func serveNews(w http.ResponseWriter, r *http.Request) {

	/*
		GallantOne user hits route:
		 0 - checks for existing files
		 1 - checks for file time
		 2 - If non-exist / time over :
		    2a - call webhose api
		    2b - save webhose json output 
		 3 - If not just return the file 
	*/

	newsBytes := art.PrepareNewsService()

	bytesWritten, err := fmt.Fprint(w, string(newsBytes))
	if err != nil {
		byteStr := rune(bytesWritten)
		log.Println("Error writing news with bytes: " + string(byteStr))
	}
	log.Println("Completion log from news: " + r.Method)
}


// main is main
func main() {
	fmt.Println("GallantOne Core Amd64")
	handleRequests()
}

// favicon for chi sever recognition
func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../<yourfaviconlocation>/favicon.ico")
}
