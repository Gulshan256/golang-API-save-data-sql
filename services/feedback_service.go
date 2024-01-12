// services/feedback_service.go
package services

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang-API-save-data-sql/models"
)

// FeedbackStoreData stores feedback data in the SQLite database.
func FeedbackStoreData(db *sql.DB, data models.Feedback) error {
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

// FeedbackGetData retrieves feedback data from the SQLite database.
func FeedbackGetData(db *sql.DB) ([]models.Feedback, error) {
	// Retrieve feedback data from the feedback table
	rows, err := db.Query(`
		SELECT Date, Time, Feedback
		FROM feedback
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbackList []models.Feedback

	// Iterate over the rows and scan the data into a slice of Feedback
	for rows.Next() {
		var feedbackData models.Feedback
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

// PostFeedback handles the HTTP POST request for feedback data.
func PostFeedback(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Decode the JSON payload
	var data models.Feedback
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
	err = FeedbackStoreData(db, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	response := map[string]string{"message": "Feedback received and stored successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFeedback handles the HTTP GET request for feedback data.
func GetFeedback(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Retrieve data from the SQLite database
	feedbackList, err := FeedbackGetData(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the retrieved data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbackList)
}
