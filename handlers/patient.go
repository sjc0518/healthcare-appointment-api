package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"healthcare-appointment-api/models"
)

// PatientHandler manages patient API requests.
type PatientHandler struct {
	DB *sql.DB
}

// NewPatientHandler creates a patient handler.
func NewPatientHandler(db *sql.DB) *PatientHandler {
	return &PatientHandler{DB: db}
}

// CreatePatient handles POST /api/v1/patients.
func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var patient models.Patient

	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	patient.Name = strings.TrimSpace(patient.Name)
	patient.Email = strings.TrimSpace(patient.Email)

	if patient.Name == "" || patient.Email == "" ||
		patient.DateOfBirth == "" {
		http.Error(w, "Name, email and date of birth are required",
			http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", patient.DateOfBirth); err != nil {
		http.Error(w, "Date of birth must be YYYY-MM-DD",
			http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(
		`INSERT INTO patients
		(name, date_of_birth, email, phone_number)
		VALUES (?, ?, ?, ?)`,
		patient.Name,
		patient.DateOfBirth,
		patient.Email,
		patient.PhoneNumber,
	)

	if err != nil {
		http.Error(w, "Failed to create patient",
			http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve patient ID",
			http.StatusInternalServerError)
		return
	}

	patient.ID = id

	// Read the saved record to return the database timestamp.
	var createdAt time.Time

	err = h.DB.QueryRow(
		"SELECT created_at FROM patients WHERE id = ?",
		id,
	).Scan(&createdAt)

	if err != nil {
		http.Error(w, "Patient created, but timestamp could not be read",
			http.StatusInternalServerError)
		return
	}

	patient.CreatedAt = createdAt.UTC()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

// GetPatients handles GET /api/v1/patients.
// Supports page and limit query parameters.
func (h *PatientHandler) GetPatients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page := 1
	limit := 10

	if value := r.URL.Query().Get("page"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			http.Error(w, "Invalid page", http.StatusBadRequest)
			return
		}
		page = n
	}

	if value := r.URL.Query().Get("limit"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 || n > 100 {
			http.Error(w, "Limit must be between 1 and 100",
				http.StatusBadRequest)
			return
		}
		limit = n
	}

	offset := (page - 1) * limit

	rows, err := h.DB.Query(
		`SELECT id, name, date_of_birth, email,
		phone_number, created_at
		FROM patients
		ORDER BY id
		LIMIT ? OFFSET ?`,
		limit, offset,
	)

	if err != nil {
		http.Error(w, "Failed to fetch patients",
			http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	patients := make([]models.Patient, 0)

	for rows.Next() {
		var patient models.Patient
		var dob time.Time
		var createdAt time.Time

		err := rows.Scan(
			&patient.ID,
			&patient.Name,
			&dob,
			&patient.Email,
			&patient.PhoneNumber,
			&createdAt,
		)

		if err != nil {
			http.Error(w, "Failed to read patient",
				http.StatusInternalServerError)
			return
		}

		patient.DateOfBirth = dob.Format("2006-01-02")
		patient.CreatedAt = createdAt.UTC()

		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read patients",
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(patients)
}

// GetPatientByID handles GET /api/v1/patients/{id}.
func (h *PatientHandler) GetPatientByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	var patient models.Patient
	var dob time.Time
	var createdAt time.Time

	err = h.DB.QueryRow(
		`SELECT id, name, date_of_birth, email,
		phone_number, created_at
		FROM patients
		WHERE id = ?`,
		id,
	).Scan(
		&patient.ID,
		&patient.Name,
		&dob,
		&patient.Email,
		&patient.PhoneNumber,
		&createdAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Database error",
			http.StatusInternalServerError)
		return
	}

	patient.DateOfBirth = dob.Format("2006-01-02")
	patient.CreatedAt = createdAt.UTC()

	json.NewEncoder(w).Encode(patient)
}
