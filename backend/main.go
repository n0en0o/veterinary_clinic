package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"my-docker-app/backend/graphics"
	"my-docker-app/backend/models"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//ждемс бдшку
func (a *App) waitForDB(host, user, password, dbname string) error {
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)

	var err error
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		a.DB, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Printf("Failed to open database: %v. Retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = a.DB.Ping()
		if err != nil {
			log.Printf("Database not ready: %v. Retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Successfully connected to database!")
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts", maxAttempts)
}



func (a *App) Initialize(host, user, password, dbname string) {

	err := a.waitForDB(host, user, password, dbname)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}



func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/health", a.healthCheck).Methods("GET")
	
	
	a.Router.HandleFunc("/api/owners", a.getOwners).Methods("GET")
	a.Router.HandleFunc("/api/owners/{id}", a.getOwner).Methods("GET")
	a.Router.HandleFunc("/api/owners", a.createOwner).Methods("POST")
	
	
	a.Router.HandleFunc("/api/pets", a.getPets).Methods("GET")
	a.Router.HandleFunc("/api/pets/{id}", a.getPet).Methods("GET")
	a.Router.HandleFunc("/api/pets/owner/{ownerId}", a.getPetsByOwner).Methods("GET")
	a.Router.HandleFunc("/api/pets", a.createPet).Methods("POST")
	
	
	a.Router.HandleFunc("/api/health-records/pet/{petId}", a.getHealthRecords).Methods("GET")
	a.Router.HandleFunc("/api/health-records", a.createHealthRecord).Methods("POST")
	
	a.Router.HandleFunc("/api/pets/{petId:[0-9]+}/chart", a.getPetChart).Methods("GET")

	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("/app/static/")))
}


func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	var dbStatus string
	err := a.DB.Ping()
	if err != nil {
		dbStatus = "disconnected"
	} else {
		dbStatus = "connected"
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"database":  dbStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func (a *App) getOwners(w http.ResponseWriter, r *http.Request) {
	owners, err := getOwners(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, owners)
}

func (a *App) getOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}

	owner, err := getOwner(a.DB, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Owner not found")
		return
	}
	respondWithJSON(w, http.StatusOK, owner)
}

func (a *App) createOwner(w http.ResponseWriter, r *http.Request) {
	var owner models.Owner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&owner); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := owner.CreateOwner(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, owner)
}


func (a *App) getPets(w http.ResponseWriter, r *http.Request) {
	pets, err := getPets(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, pets)
}

func (a *App) getPet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	pet, err := getPet(a.DB, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Pet not found")
		return
	}
	respondWithJSON(w, http.StatusOK, pet)
}

func (a *App) getPetsByOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerId, err := strconv.Atoi(vars["ownerId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}

	pets, err := getPetsByOwner(a.DB, ownerId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, pets)
}

func (a *App) createPet(w http.ResponseWriter, r *http.Request) {
	var pet models.Pet
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := pet.CreatePet(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, pet)
}


func (a *App) getHealthRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	petId, err := strconv.Atoi(vars["petId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	records, err := getHealthRecords(a.DB, petId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, records)
}

func (a *App) createHealthRecord(w http.ResponseWriter, r *http.Request) {
	var record models.HealthRecord
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&record); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := record.CreateHealthRecord(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, record)
}


func getOwners(db *sql.DB) ([]models.Owner, error) {
	rows, err := db.Query("SELECT id, first_name, last_name, email, phone, address, created_at FROM owners")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owners []models.Owner
	for rows.Next() {
		var owner models.Owner
		if err := rows.Scan(&owner.ID, &owner.FirstName, &owner.LastName, &owner.Email, &owner.Phone, &owner.Address, &owner.CreatedAt); err != nil {
			return nil, err
		}
		owners = append(owners, owner)
	}
	return owners, nil
}

func getOwner(db *sql.DB, id int) (models.Owner, error) {
	var owner models.Owner
	err := db.QueryRow("SELECT id, first_name, last_name, email, phone, address, created_at FROM owners WHERE id = $1", id).
		Scan(&owner.ID, &owner.FirstName, &owner.LastName, &owner.Email, &owner.Phone, &owner.Address, &owner.CreatedAt)
	return owner, err
}

func getPets(db *sql.DB) ([]models.Pet, error) {
	rows, err := db.Query(`
		SELECT p.id, p.owner_id, p.name, p.species, p.breed, p.date_of_birth, p.color, p.microchip_id, p.created_at,
		       o.first_name || ' ' || o.last_name as owner_name
		FROM pets p
		JOIN owners o ON p.owner_id = o.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []models.Pet
	for rows.Next() {
		var pet models.Pet
		if err := rows.Scan(&pet.ID, &pet.OwnerID, &pet.Name, &pet.Species, &pet.Breed, &pet.DateOfBirth, &pet.Color, &pet.MicrochipID, &pet.CreatedAt, &pet.OwnerName); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return pets, nil
}

func getPet(db *sql.DB, id int) (models.Pet, error) {
	var pet models.Pet
	err := db.QueryRow(`
		SELECT p.id, p.owner_id, p.name, p.species, p.breed, p.date_of_birth, p.color, p.microchip_id, p.created_at,
		       o.first_name || ' ' || o.last_name as owner_name
		FROM pets p
		JOIN owners o ON p.owner_id = o.id
		WHERE p.id = $1
	`, id).Scan(&pet.ID, &pet.OwnerID, &pet.Name, &pet.Species, &pet.Breed, &pet.DateOfBirth, &pet.Color, &pet.MicrochipID, &pet.CreatedAt, &pet.OwnerName)
	return pet, err
}

func getPetsByOwner(db *sql.DB, ownerId int) ([]models.Pet, error) {
	rows, err := db.Query(`
		SELECT id, owner_id, name, species, breed, date_of_birth, color, microchip_id, created_at
		FROM pets WHERE owner_id = $1
	`, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []models.Pet
	for rows.Next() {
		var pet models.Pet
		if err := rows.Scan(&pet.ID, &pet.OwnerID, &pet.Name, &pet.Species, &pet.Breed, &pet.DateOfBirth, &pet.Color, &pet.MicrochipID, &pet.CreatedAt); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return pets, nil
}


func getHealthRecords(db *sql.DB, petId int) ([]models.HealthRecord, error) {
	rows, err := db.Query(`
		SELECT id, pet_id, visit_date, weight, temperature, heart_rate, respiratory_rate, notes, diagnosis, treatment, next_visit_date, created_at
		FROM health_records WHERE pet_id = $1 ORDER BY visit_date DESC
	`, petId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.HealthRecord
	for rows.Next() {
		var record models.HealthRecord
		if err := rows.Scan(&record.ID, &record.PetID, &record.VisitDate, &record.Weight, &record.Temperature, &record.HeartRate, &record.RespiratoryRate, &record.Notes, &record.Diagnosis, &record.Treatment, &record.NextVisitDate, &record.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}


func (a *App) getPetChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	petId, err := strconv.Atoi(vars["petId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	records, err := getHealthRecords(a.DB, petId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(records) == 0 {
        respondWithError(w, http.StatusNotFound, "Для этого питомца еще нет записей о здоровье")
        return
    }

	w.Header().Set("Content-Type", "text/html")

	graphics.DrawPetHealthChart(records, w)

}

//вспомогательные функции
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//запуск приложения
func (a *App) Run(addr string) {
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func main() {
	a := &App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	a.Run(":8000")
}