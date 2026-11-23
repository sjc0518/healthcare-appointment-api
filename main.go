package main

import (
	"encoding/json"
	"log"
	"net/http"

	"healthcare-appointment-api/config"
	"healthcare-appointment-api/routes"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	defer db.Close()

	log.Println("MySQL database connected successfully!")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "Database unavailable",
				http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{
			"status":   "UP",
			"database": "CONNECTED",
			"message":  "Healthcare Appointment API is running",
		})
	})

	routes.RegisterRoutes(mux, db)

	log.Println("Server starting on port 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
