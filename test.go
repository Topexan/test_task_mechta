package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	//"math/rand"
	"net/http"
	//"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
	Code         string `json:"code"`
	Country_code string `json:"country_code"`
}

func listCities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cities []City
	db := setupDB()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT id, name, code, country_code FROM cities")
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var code string
		var country_code string
		err = rows.Scan(&id, &name, &code, &country_code)
		if err != nil {
			// handle this error
			panic(err)
		}
		cities = append(cities, City{ID: id, Name: name, Code: code, Country_code: country_code})
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	log.Println("Get all cities from table: cities")
	json.NewEncoder(w).Encode(cities)
}

//function to show one person
func getCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //read parameter from url

	db := setupDB()
	defer db.Close()

	sqlStatement := `
		SELECT id, name, code, country_code FROM cities
		WHERE id = $1;`
	rows, err := db.Query(sqlStatement, params["id"])
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	var c City
	for rows.Next() {
		var id int
		var name string
		var code string
		var country_code string
		err = rows.Scan(&id, &name, &code, &country_code)
		if err != nil {
			// handle this error
			panic(err)
		}
		c = City{ID: id, Name: name, Code: code, Country_code: country_code}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	log.Println("Get city from table: cities, with id:", params["id"])
	json.NewEncoder(w).Encode(c)
}

func createCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var c City //new City
	err := decoder.Decode(&c)

	if err != nil {
		panic(err)
	}

	_ = json.NewDecoder(r.Body).Decode(&c)

	//p.ID = strconv.Itoa(rand.Intn(30-4) + 4) // random number from 4 to 30
	json.NewEncoder(w).Encode(c)

	db := setupDB()
	defer db.Close()

	sqlStatement := `
		INSERT INTO cities (name, code, country_code)
		VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, c.Name, c.Code, c.Country_code) //add city to database
	if err != nil {
		panic(err)
	}
	log.Println("Insert new city to table: cities")
}

// function to update existing person
func updateCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var c City //person with new parameters
	err := decoder.Decode(&c)

	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(c)

	db := setupDB()
	defer db.Close()

	sqlStatement := `
		UPDATE cities SET name = $2,
						   code = $3,
						   country_code = $4
		WHERE id = $1;`
	_, err = db.Exec(sqlStatement, c.ID, c.Name, c.Code, c.Country_code) //update city in our database
	if err != nil {
		panic(err)
	}
	log.Println("Update person from table: persons, with id:", c.ID)

}

func deleteCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //take id from url

	db := setupDB()
	defer db.Close()

	sqlStatement := `
		DELETE FROM cities
		WHERE id = $1;`
	_, err := db.Exec(sqlStatement, params["id"]) //delete city from database
	if err != nil {
		panic(err)
	}
	log.Println("Delete city from table: cities, with id:", params["id"])
	listCities(w, r)
}

func setupDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		//panic(err)
		log.Println(err)
	}
	log.Println("Succesful connetcion to ", dbname)
	return db
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
