package routes

import (
	"database/sql"
	"net/http"

	"healthcare-appointment-api/handlers"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {

	// Initialize Patient Handler
	patientHandler := &handlers.PatientHandler{
		DB: db,
	}

	// POST - Create a new patient
	mux.HandleFunc(
		"POST /api/v1/patients",
		patientHandler.CreatePatient,
	)

	// GET - List all patients with pagination
	mux.HandleFunc(
		"GET /api/v1/patients",
		patientHandler.GetPatients,
	)

	// GET - Retrieve a patient by ID
	mux.HandleFunc(
		"GET /api/v1/patients/{id}",
		patientHandler.GetPatientByID,
	)
}
