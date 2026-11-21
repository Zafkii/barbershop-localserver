package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"barbershop-localserver/db"
	"barbershop-localserver/handlers"
)

func main() {
	path := os.Getenv("BARBERHOP_DB")
	if path == "" {
		path = "./barbershop.db"
	}

	database, err := db.Open(path)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	router := handlers.NewRouter(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Printf("localserver listening on %s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
