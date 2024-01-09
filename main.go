package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Data represents the structure of the data to be stored.
type Data struct {
	Date          string `json:"date"`
	Time          string `json:"time"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	FromLocation  string `json:"from"`
	ComponentName string `json:"componentName"`
}

var db *sql.DB

func main() {
	fmt.Println("Starting server...")
	// Open the SQLite database
	var err error
	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Create the data table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS data (
			Date TEXT,
			Time TEXT,
			Name TEXT,
			Email TEXT,
			Phone TEXT,
			FromLocation TEXT,
			ComponentName TEXT
		)
	`)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()

	// Define a route for handling POST requests to /post-data
	r.HandleFunc("/post-data", postDataHandler).Methods("POST")

	// Define a route for handling GET requests to /get-data
	r.HandleFunc("/get-data", getDataHandler).Methods("GET")

	// Start the server on port 8080
	http.Handle("/", r)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func postDataHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON payload
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if data.Name == "" || data.Email == "" || data.Phone == "" || data.FromLocation == "" || data.ComponentName == "" {
		http.Error(w, "Required fields cannot be empty", http.StatusBadRequest)
		return
	}

	// Set Date and Time to current date and time
	data.Date = time.Now().Format("2006-01-02")
	data.Time = time.Now().Format("15:04:05")

	// Store the data in the SQLite database
	err = storeData(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	response := map[string]string{"message": "Data received and stored successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve data from the SQLite database
	dataList, err := getData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the retrieved data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataList)
}

func storeData(data Data) error {
	// Prepare the SQL statement for inserting data
	stmt, err := db.Prepare(`
		INSERT INTO data (Date, Time, Name, Email, Phone, FromLocation, ComponentName)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with data values
	_, err = stmt.Exec(data.Date, data.Time, data.Name, data.Email, data.Phone, data.FromLocation, data.ComponentName)
	if err != nil {
		return err
	}

	return nil
}

func getData() ([]Data, error) {
	// Retrieve all data from the data table
	rows, err := db.Query(`
		SELECT Date, Time, Name, Email, Phone, FromLocation, ComponentName
		FROM data
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []Data

	// Iterate over the rows and scan the data into a slice of Data
	for rows.Next() {
		var data Data
		err := rows.Scan(
			&data.Date,
			&data.Time,
			&data.Name,
			&data.Email,
			&data.Phone,
			&data.FromLocation,
			&data.ComponentName,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}
