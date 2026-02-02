package models

import (
	"database/sql"
	"time"
)

type Owner struct {
    ID        int       `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Address   string    `json:"address"`
    CreatedAt time.Time `json:"created_at"`
}

type Pet struct {
    ID          int       `json:"id"`
    OwnerID     int       `json:"owner_id"`
    Name        string    `json:"name"`
    Species     string    `json:"species"`
    Breed       string    `json:"breed"`
    DateOfBirth string    `json:"date_of_birth"`
    Color       string    `json:"color"`
    MicrochipID string    `json:"microchip_id"`
    CreatedAt   time.Time `json:"created_at"`
    OwnerName   string    `json:"owner_name,omitempty"`
}

type HealthRecord struct {
    ID              int       `json:"id"`
    PetID           int       `json:"pet_id"`
    VisitDate       string    `json:"visit_date"`
    Weight          float64   `json:"weight"`
    Temperature     float64   `json:"temperature"`
    HeartRate       int       `json:"heart_rate"`
    RespiratoryRate int       `json:"respiratory_rate"`
    Notes           string    `json:"notes"`
    Diagnosis       string    `json:"diagnosis"`
    Treatment       string    `json:"treatment"`
    NextVisitDate   string    `json:"next_visit_date"`
    CreatedAt       time.Time `json:"created_at"`
}

func (owner *Owner) CreateOwner(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO owners(first_name, last_name, email, phone, address) VALUES($1, $2, $3, $4, $5) RETURNING id",
		owner.FirstName, owner.LastName, owner.Email, owner.Phone, owner.Address,
	).Scan(&owner.ID)
	return err
}

func (pet *Pet) CreatePet(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO pets(owner_id, name, species, breed, date_of_birth, color, microchip_id) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		pet.OwnerID, pet.Name, pet.Species, pet.Breed, pet.DateOfBirth, pet.Color, pet.MicrochipID,
	).Scan(&pet.ID)
	return err
}

func (record *HealthRecord) CreateHealthRecord(db *sql.DB) error { 
	err := db.QueryRow(
		"INSERT INTO health_records(pet_id, visit_date, weight, temperature, heart_rate, respiratory_rate, notes, diagnosis, treatment, next_visit_date) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		record.PetID, record.VisitDate, record.Weight, record.Temperature, record.HeartRate, record.RespiratoryRate, record.Notes, record.Diagnosis, record.Treatment, record.NextVisitDate,
	).Scan(&record.ID)
	return err
}

