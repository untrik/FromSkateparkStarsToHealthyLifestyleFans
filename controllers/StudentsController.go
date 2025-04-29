package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/untrik/FromSkateToZOH/database"
	"github.com/untrik/FromSkateToZOH/middleware"
	"github.com/untrik/FromSkateToZOH/models"
)

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Faculty  string `json:"faculty"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if request.Username == "" || request.Password == "" || request.Name == "" || request.LastName == "" || request.Faculty == "" {
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
	user := models.User{
		Username:     request.Username,
		PasswordHash: hashedPassword,
	}
	var existingUser models.User
	if err := database.DB.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		log.Print("Username already exists")
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}
	if err = database.DB.Create(&user).Error; err != nil {
		log.Print("database creation user error", err)
		http.Error(w, "database creation user error", http.StatusBadRequest)
		return
	}
	student := models.Student{
		UserID:   user.ID,
		Name:     request.Name,
		LastName: request.LastName,
		Faculty:  request.Faculty,
	}
	if err = database.DB.Create(&student).Error; err != nil {
		log.Print("database creation student error", err)
		http.Error(w, "database creation student error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}
func RegistrationForTheEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	IdStudent, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		log.Print("conversion error")
		http.Error(w, "conversion errorr", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	IdEvent, err := strconv.Atoi(vars["id_event"])
	if err != nil {
		log.Print("conversion error")
		http.Error(w, "conversion errorr", http.StatusBadRequest)
		return
	}
	var student models.Student
	if err = database.DB.Where("user_id = ?", IdStudent).First(&student).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	var event models.Event
	if err = database.DB.Where("event_id = ?", IdEvent).First(&event).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}

	eventParticipant := models.EventParticipant{
		EventID:   uint(IdEvent),
		StudentID: student.StudentId,
	}
	if err = database.DB.Create(&eventParticipant).Error; err != nil {
		log.Print("database registration student error", err)
		http.Error(w, "database registration student error", http.StatusBadRequest)
		return
	}
	type RegistrationResponse struct {
		EventID   uint                 `json:"event_id"`
		StudentID uint                 `json:"student_id"`
		Status    models.StatusStudent `json:"status"`
	}
	response := RegistrationResponse{
		EventID:   eventParticipant.EventID,
		StudentID: eventParticipant.StudentID,
		Status:    eventParticipant.Status,
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}
