package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "enteam"
	dbname   = "test_task_mechta"
)

type City struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"surname"`
	Country_code string `json:"city"`
}

func main() {
	r := mux.NewRouter() //create new router
	r.HandleFunc("/cities", listCities).Methods("GET")
	r.HandleFunc("/cities/{id}", getCity).Methods("GET")
	r.HandleFunc("/cities", createCity).Methods("POST")
	r.HandleFunc("/cities/{id}", updateCity).Methods("PUT")
	r.HandleFunc("/cities/{id}", deleteCity).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
