package models

import "time"

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