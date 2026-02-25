package main

import (
	"fmt"
	"log"
	"net/http"

	"myapp/internal/config"
	"myapp/internal/handlers"
	"myapp/internal/store"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.GetConfig()

	db, err := store.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()

	h := handlers.NewHandler(db)

	router.HandleFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/users", h.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

	addr := ":8080"
	fmt.Printf("Server is running at %s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
