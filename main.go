package main

import (
	_ "fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"time"
)

func main() {
	router := mux.NewRouter()
	setupRoutes(router)

	port := ":8000"
	server := &http.Server{
		Handler: router,
		Addr:    port,
		// SET TIMEOUT
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server started at %s", port)
	log.Fatal(server.ListenAndServe())
}