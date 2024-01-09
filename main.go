package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// Data represents the structure of the data to be stored.
type Waitlist struct {
	Date          string `json:"date"`
	Time          string `json:"time"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	FromLocation  string `json:"from"`
	ComponentName string `json:"componentName"`
}
type Feedback struct {
	Date     string `json:"date"`
	Time     string `json:"time"`
	Feedback string `json:"feedback"`
}


var db *sql.DB



func main() {


	loadConfig()

	// Open the SQLite database
	var err error
	db, err = sql.Open(viper.GetString("database.driver"), viper.GetString("database.connection"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()


	// Create the waitlist table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS waitlist (
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

	// Create the feedback table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS feedback (
			Date TEXT,
			Time TEXT,
			Feedback TEXT
		)
	`)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()

	// Define a route for handling POST requests to /post-data
	r.HandleFunc("/post-waitlist-data", postWaitlistData).Methods("POST")

	// Define a route for handling GET requests to /get-waitlist-data
	r.HandleFunc("/get-waitlist-data", getWaitlistData).Methods("GET")

	// Define a route for handling POST requests to /post-feedback
	r.HandleFunc("/post-feedback", postFeedback).Methods("POST")

	// Define a route for handling GET requests to /get-feedback
	r.HandleFunc("/get-feedback", getFeedback).Methods("GET")

	// Start the server on the specified port
	port := viper.GetString("server.port")
	http.Handle("/", r)
	// fmt.Printf("Server listening on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}



func loadConfig() {
	// Set the file name of the configuration file
	viper.SetConfigFile("config.yaml")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}
}




func postWaitlistData(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON payload
	var data Waitlist
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
	err = WaitlistStoreData(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	response := map[string]string{"message": "Data received and stored successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}



func getWaitlistData(w http.ResponseWriter, r *http.Request) {
	// Retrieve data from the SQLite database
	dataList, err := WaitlistGetData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the retrieved data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataList)
}




func WaitlistStoreData(data Waitlist) error {
	// Prepare the SQL statement for inserting data
	stmt, err := db.Prepare(`
		INSERT INTO waitlist (Date, Time, Name, Email, Phone, FromLocation, ComponentName)
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




func WaitlistGetData() ([]Waitlist, error) {
	// Retrieve all data from the data table
	rows, err := db.Query(`
		SELECT Date, Time, Name, Email, Phone, FromLocation, ComponentName
		FROM waitlist
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []Waitlist

	// Iterate over the rows and scan the data into a slice of Data
	for rows.Next() {
		var data Waitlist
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




func postFeedback(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON payload
	var data Feedback
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if data.Feedback == "" {
		http.Error(w, "Feedback cannot be empty", http.StatusBadRequest)
		return
	}

	// Set Date and Time to current date and time
	data.Date = time.Now().Format("2006-01-02")
	data.Time = time.Now().Format("15:04:05")

	// Store the data in the SQLite database
	err = FeedbackStoreData(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	response := map[string]string{"message": "Feedback received and stored successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}




func getFeedback(w http.ResponseWriter, r *http.Request) {
	// Retrieve data from the SQLite database
	feedbackList, err := FeedbackGetData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the retrieved data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbackList)
}





func FeedbackStoreData(data Feedback) error {
	// Prepare the SQL statement for inserting data
	stmt, err := db.Prepare(`
		INSERT INTO feedback (Date, Time, Feedback)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with data values
	_, err = stmt.Exec(data.Date, data.Time, data.Feedback)
	if err != nil {
		return err
	}

	return nil
}



func FeedbackGetData() ([]Feedback, error) {
	// Retrieve feedback data from the feedback table
	rows, err := db.Query(`
		SELECT Date, Time, Feedback
		FROM feedback
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbackList []Feedback

	// Iterate over the rows and scan the data into a slice of Data
	for rows.Next() {
		var feedbackData Feedback
		err := rows.Scan(
			&feedbackData.Date,
			&feedbackData.Time,
			&feedbackData.Feedback,
		)
		if err != nil {
			return nil, err
		}

		feedbackList = append(feedbackList, feedbackData)
	}

	return feedbackList, nil
}
