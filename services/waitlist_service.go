// services/waitlist_service.go
package services

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang-API-save-data-sql/models"
)

// WaitlistStoreData stores waitlist data in the SQLite database.
func WaitlistStoreData(db *sql.DB, data models.Waitlist) error {
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

// WaitlistGetData retrieves waitlist data from the SQLite database.
func WaitlistGetData(db *sql.DB) ([]models.Waitlist, error) {
	// Retrieve all data from the waitlist table
	rows, err := db.Query(`
		SELECT Date, Time, Name, Email, Phone, FromLocation, ComponentName
		FROM waitlist
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []models.Waitlist

	// Iterate over the rows and scan the data into a slice of Waitlist
	for rows.Next() {
		var data models.Waitlist
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

// PostWaitlistData handles the HTTP POST request for waitlist data.
func PostWaitlistData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Decode the JSON payload
	var data models.Waitlist
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
	err = WaitlistStoreData(db, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	response := map[string]string{"message": "Data received and stored successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetWaitlistData handles the HTTP GET request for waitlist data.
func GetWaitlistData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Retrieve data from the SQLite database
	dataList, err := WaitlistGetData(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the retrieved data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataList)
}
