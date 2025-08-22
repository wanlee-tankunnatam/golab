package handler

import (
	"encoding/json"
	"log"
	"net/http"

	config "go-lab/config"
	model "go-lab/model"

	faker "github.com/go-faker/faker/v4"
)

func InsertUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user := model.Users{
		Name:  faker.Name(),
		Email: faker.Email(),
	}

	log.Printf("Inserting user: %s %s", user.Name, user.Email)
	if result := config.DB.Create(&user); result.Error != nil {
		log.Printf("DB error: %v", result.Error)
		http.Error(w, "Database Error: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
