package main

import (
	routers "golang-API-save-data-sql/routes"
	"log"

	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var DB *sql.DB // Exported variable

func main() {
	loadConfig()

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Open the SQLite database
	var err error
	DB, err = sql.Open(viper.GetString("database.driver"), viper.GetString("database.connection"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer DB.Close()

	// Create the waitlist table if it doesn't exist
	_, err = DB.Exec(`
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
		log.Fatal(err)
		return
	}

	// Create the feedback table if it doesn't exist
	_, err = DB.Exec(`
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

	// Setup routes
	routers.SetupRoutes(r, DB)

	// Start the server on the specified port in .env file
	envFile, _ := godotenv.Read(".env")
	port := envFile["port"]

	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}

func loadConfig() {
	// Set the file name of the configuration file
	viper.SetConfigFile("config.yaml")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}
}
