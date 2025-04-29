package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/untrik/FromSkateToZOH/database"
	"github.com/untrik/FromSkateToZOH/middleware"
	"github.com/untrik/FromSkateToZOH/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if request.Username == "" || request.Password == "" {
		log.Print("Missing request fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := database.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	if unHashing(user.PasswordHash, request.Password) {
		token, err := middleware.GenerateJWT(user.ID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token":      token,
			"expires_in": 3600 * 12,
			"token_type": "Bearer",
		})

	} else {
		log.Print("invalid password")
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

}
func unHashing(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
func GetAllActiveEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var events []models.Event
	if err := database.DB.Where("date > ?", time.Now()).Find(&events).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}
