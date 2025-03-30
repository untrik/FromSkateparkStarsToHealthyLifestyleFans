package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/untrik/FromSkateToZOH/database"
	"github.com/untrik/FromSkateToZOH/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		models.Student
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if request.Name == "" || request.LastName == "" || request.Username == "" || request.Faculty == "" || request.Password == "" {
		log.Print("Missing request fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	hashedPassword, err := hashing(request.Password)
	if err != nil {
		log.Print("Password hashing failed: ", err)
		http.Error(w, "Password hashing failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	user := models.Student{
		Username:     request.Username,
		Name:         request.Name,
		LastName:     request.LastName,
		Faculty:      request.Faculty,
		PasswordHash: hashedPassword,
	}

	fmt.Println(user)
	if err = database.DB.Create(&user).Error; err != nil {
		log.Print("database creation user error", err)
		http.Error(w, "database creation user error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}
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
	var user models.Student
	if err := database.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		log.Print("No user in db", err)
		http.Error(w, "No user in db: "+err.Error(), http.StatusBadRequest)
		return
	}
	if unHashing(user.PasswordHash, request.Password) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
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
func hashing(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is nil")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}
