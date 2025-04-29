package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/untrik/FromSkateToZOH/database"
	"github.com/untrik/FromSkateToZOH/middleware"
	"github.com/untrik/FromSkateToZOH/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	createdBy, ok := r.Context().Value(middleware.UserIDKey).(uint)
	fmt.Println(createdBy)
	if !ok {
		log.Print("Missing ID parameter")
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	var request struct {
		Username string
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
	hashedPassword, err := hashing(request.Password)
	if err != nil {
		log.Print("Password hashing failed: ", err)
		http.Error(w, "Password hashing failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	var admin models.Admin

	if err := database.DB.Where("user_id", createdBy).First(&admin).Error; err != nil {
		log.Print("The creator's admin does not exist")
		http.Error(w, "The creator's admin does not exist", http.StatusBadRequest)
		return
	}
	user := models.User{
		Username:     request.Username,
		PasswordHash: hashedPassword,
	}
	if err = database.DB.Create(&user).Error; err != nil {
		log.Print("database creation user error", err)
		http.Error(w, "database creation user error", http.StatusBadRequest)
		return
	}

	admin = models.Admin{
		UserID:    user.ID,
		CreatedBy: uint(createdBy),
	}

	if err = database.DB.Create(&admin).Error; err != nil {
		log.Print("database creation user error", err)
		http.Error(w, "database creation user error", http.StatusBadRequest)
		return
	}
	type AdminResponse struct {
		AdminID   uint        `json:"admin_id"`
		UserID    uint        `json:"user_id"`
		User      models.User `json:"user"`
		CreatedBy uint        `json:"created_by"`
	}
	adminResponse := AdminResponse{
		AdminID:   admin.AdminId,
		UserID:    user.ID,
		User:      user,
		CreatedBy: uint(createdBy),
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(adminResponse)
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
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	IdUser, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		log.Print("conversion error")
		http.Error(w, "conversion error", http.StatusBadRequest)
		return
	}
	var request struct {
		Title       string `json:"title"`
		Location    string `json:"location"`
		Description string `json:"description"`
		Date        string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if request.Date == "" || request.Description == "" || request.Location == "" || request.Title == "" {
		log.Print("Missing request fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	var admin models.Admin
	if err := database.DB.Where("user_id", IdUser).First(&admin).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	layout := "2006-01-02 15:04:05"
	data, err := time.Parse(layout, request.Date)
	if err != nil {

	}
	event := models.Event{
		Title:       request.Title,
		Date:        data,
		Location:    request.Location,
		Description: request.Description,
		AdminId:     admin.AdminId,
	}
	if err = database.DB.Create(&event).Error; err != nil {
		log.Print("database creation event error", err)
		http.Error(w, "database creation event error", http.StatusBadRequest)
		return
	}
	type eventResponse struct {
		EventId uint   `json:"event_id"`
		Title   string `json:"title"`

		Date        time.Time `json:"date"`
		Location    string    `json:"location"`
		Description string    `json:"description"`
	}
	response := eventResponse{
		EventId:     event.EventId,
		Title:       event.Title,
		Date:        event.Date,
		Location:    event.Location,
		Description: event.Description,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func AddReward(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	IdStudent, err := strconv.Atoi(vars["id_student"])
	if err != nil {
		log.Print("conversion error")
		http.Error(w, "conversion error", http.StatusBadRequest)
		return
	}
	IdEvent, err := strconv.Atoi(vars["id_event"])
	if err != nil {
		log.Print("conversion error")
		http.Error(w, "conversion error", http.StatusBadRequest)
		return
	}
	var student models.Student
	if err = database.DB.Where("student_id = ?", IdStudent).First(&student).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	var eventParticipant models.EventParticipant
	if err = database.DB.Where("student_id = ?", student.StudentId).Where("event_id", IdEvent).First(&eventParticipant).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	if eventParticipant.Status != "approved" {
		log.Print("Student status not approved")
		http.Error(w, "The student did not attend the event", http.StatusBadRequest)
		return
	}

	var event models.Event
	if err = database.DB.Where("event_id = ?", IdEvent).First(&event).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	var request struct {
		Place  int `json:"place"`
		Points int `json:"points"`
	}
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if request.Place <= 0 || request.Points <= 0 {
		log.Print("Missing request fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	reward := models.Reward{
		EventID:   uint(IdEvent),
		Place:     request.Place,
		Points:    request.Points,
		StudentID: uint(IdStudent),
	}
	if err = database.DB.Create(&reward).Error; err != nil {
		log.Print("database create reward error", err)
		http.Error(w, "database create reward error", http.StatusBadRequest)
		return
	}
	points := student.Points + uint(request.Points)
	database.DB.Model(&student).Where("student_id", IdStudent).Update("points", points)
	type responseReward struct {
		EventID   uint `json:"event_id"`
		Place     int  `json:"place"`
		Points    int  `json:"points"`
		StudentID uint `json:"student_id"`
	}
	response := responseReward{
		EventID:   reward.EventID,
		Place:     request.Place,
		Points:    request.Points,
		StudentID: reward.StudentID,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func ChangeStatusStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		Status models.StatusStudent `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	switch request.Status {
	case models.StatusRegistered, models.StatusApproved, models.StatusCancelled:

	default:
		log.Print("Invalid status")
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	vars := mux.Vars(r)
	StudentId, err := strconv.Atoi(vars["id_student"])
	if err != nil {
		log.Print("conversion error", err)
		http.Error(w, "conversion error"+err.Error(), http.StatusBadRequest)
		return
	}
	EventID, err := strconv.Atoi(vars["id_event"])
	if err != nil {
		log.Print("conversion error", err)
		http.Error(w, "conversion error"+err.Error(), http.StatusBadRequest)
		return
	}
	var eventParticipant models.EventParticipant
	if err := database.DB.Model(&eventParticipant).Where("student_id = ?", StudentId).Where("event_id = ?", EventID).Update("status", request.Status).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	eventParticipant.EventID = uint(EventID)
	eventParticipant.StudentID = uint(StudentId)
	if err := database.DB.
		Preload("Student").
		Preload("Event").
		Where("student_id = ? AND event_id = ?", StudentId, EventID).
		First(&eventParticipant).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	type ResponseStatus struct {
		Status  models.StatusStudent
		Student struct {
			StudentId uint   `json:"id"`
			Name      string `json:"name"`
			LastName  string `json:"last_name"`
			Faculty   string `json:"faculty"`
		}
		Event struct {
			EventID     uint      `json:"event_id"`
			Title       string    `json:"title"`
			Location    string    `json:"location"`
			Description string    `json:"description"`
			Date        time.Time `json:"date"`
		}
	}
	response := ResponseStatus{
		Status: request.Status,
		Student: struct {
			StudentId uint   `json:"id"`
			Name      string `json:"name"`
			LastName  string `json:"last_name"`
			Faculty   string `json:"faculty"`
		}{
			StudentId: eventParticipant.StudentID,
			Name:      eventParticipant.Student.Name,
			LastName:  eventParticipant.Student.LastName,
			Faculty:   eventParticipant.Student.Faculty,
		},
		Event: struct {
			EventID     uint      `json:"event_id"`
			Title       string    `json:"title"`
			Location    string    `json:"location"`
			Description string    `json:"description"`
			Date        time.Time `json:"date"`
		}{
			EventID:     eventParticipant.EventID,
			Title:       eventParticipant.Event.Title,
			Location:    eventParticipant.Event.Location,
			Description: eventParticipant.Event.Description,
			Date:        eventParticipant.Event.Date,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	eventId, err := strconv.Atoi(vars["id_event"])
	if err != nil {
		log.Print("conversion error", err)
		http.Error(w, "conversion error"+err.Error(), http.StatusBadRequest)
		return
	}
	var event models.Event
	if err := database.DB.Where("event_id = ? AND date > ?", eventId, time.Now()).Delete(&event).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	var eventParticipants models.EventParticipant
	if err := database.DB.Where("event_id = ?", eventId).Delete(&eventParticipants).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func AddProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		Price       uint   `json:"price"`
		Name        string `json:"name" `
		Description string `json:"description"`
		ImageURL    string `json:"image_URL"`
		Quantity    uint   `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if request.ImageURL == "" || request.Name == "" || request.Price <= 0 {
		log.Print("Missing request fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	product := models.Product{
		Price:       request.Price,
		Name:        request.Name,
		Description: request.Description,
		ImageURL:    request.ImageURL,
		Quantity:    request.Quantity,
	}
	rez := database.DB.Create(&product)
	if rez.Error != nil {
		log.Print("Create product Error")
		http.Error(w, "Create product Error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
func GetAllParticipants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var eventParticipants []models.EventParticipant
	idEvent, err := strconv.Atoi(mux.Vars(r)["id_event"])
	if err != nil {
		log.Print("conversion error", err)
		http.Error(w, "conversion error"+err.Error(), http.StatusBadRequest)
		return
	}
	if err := database.DB.Preload("Student").Where("event_id = ?", idEvent).Find(&eventParticipants).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	type responseParticipants struct {
		EventID uint `json:"id"`
		Student struct {
			StudentId uint                 `json:"student_id"`
			Name      string               `json:"name"`
			LastName  string               `json:"last_name"`
			Status    models.StatusStudent `json:"status"`
			Faculty   string               `json:"faculty"`
		} `json:"student"`
	}
	var response []responseParticipants
	for _, participants := range eventParticipants {
		response = append(response, responseParticipants{
			EventID: participants.EventID,
			Student: struct {
				StudentId uint                 `json:"student_id"`
				Name      string               `json:"name"`
				LastName  string               `json:"last_name"`
				Status    models.StatusStudent `json:"status"`
				Faculty   string               `json:"faculty"`
			}{
				StudentId: participants.StudentID,
				Name:      participants.Student.Name,
				LastName:  participants.Student.LastName,
				Status:    participants.Status,
				Faculty:   participants.Student.Faculty,
			},
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
