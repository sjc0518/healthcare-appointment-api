package models

import "time"

type Patient struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth string    `json:"date_of_birth"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}
