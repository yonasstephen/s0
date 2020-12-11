package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/yonasstephen/s0/store/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// setup static directory
	fs := http.FileServer(http.Dir("./static"))

	// setup router
	r := mux.NewRouter()
	r.HandleFunc("/upload", handler.UploadHandler).Methods(http.MethodPost)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Running s0 on port %s...", port)
	log.Fatal(srv.ListenAndServe())
}
