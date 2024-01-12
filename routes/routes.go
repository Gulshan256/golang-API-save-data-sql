package routers

import (
	"database/sql"
	"golang-API-save-data-sql/services"

	// "encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/post-waitlist-data", func(w http.ResponseWriter, r *http.Request) {
		services.PostWaitlistData(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/get-waitlist-data", func(w http.ResponseWriter, r *http.Request) {
		services.GetWaitlistData(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/post-feedback", func(w http.ResponseWriter, r *http.Request) {
		services.PostFeedback(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/get-feedback", func(w http.ResponseWriter, r *http.Request) {
		services.GetFeedback(w, r, db)
	}).Methods("GET")
}
