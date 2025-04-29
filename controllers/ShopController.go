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

func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var products []models.Product
	if err := database.DB.Where("quantity > 0").Where("is_deleted = ?", false).Find(&products).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
func AddOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idProduct, err := strconv.Atoi(mux.Vars(r)["id_product"])
	if err != nil {
		log.Print("conversion error")
		http.Error(w, "conversion errorr", http.StatusBadRequest)
		return
	}
	var product models.Product
	if err := database.DB.Where("product_id = ?", idProduct).First(&product).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	IdUser, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		log.Print("conversion error")
		http.Error(w, "conversion errorr", http.StatusBadRequest)
		return
	}
	var student models.Student
	if err := database.DB.Where("user_id = ?", IdUser).First(&student).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	var request struct {
		Quantity uint `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if request.Quantity == 0 || request.Quantity > product.Quantity {
		log.Print("Incorrect quantity")
		http.Error(w, "Incorrect quantity: ", http.StatusBadRequest)
		return
	}
	if student.Points < request.Quantity*product.Price {
		log.Print("not enough points")
		http.Error(w, "not enough points: ", http.StatusBadRequest)
		return
	}
	newPoints := student.Points - request.Quantity*product.Price

	order := models.Order{
		StudentId: student.StudentId,
		ProductId: uint(idProduct),
		Quantity:  request.Quantity,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		log.Print("database create order error", err)
		http.Error(w, "database create order error", http.StatusBadRequest)
		return
	}
	quantity := product.Quantity - request.Quantity
	if err := database.DB.Model(&product).Where("product_id = ?", idProduct).Update("quantity", quantity).Error; err != nil {
		log.Print("error update product", err)
		http.Error(w, "error update product: "+err.Error(), http.StatusNotFound)
		return
	}
	if err := database.DB.Model(&student).Where("user_id = ?", IdUser).Update("points", newPoints).Error; err != nil {
		log.Print("error update student", err)
		http.Error(w, "error update student: "+err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)

}
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		log.Print("conversion error")
		http.Error(w, "conversion errorr", http.StatusBadRequest)
		return
	}
	var student models.Student
	if err := database.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	var orders []models.Order
	if err := database.DB.Preload("Product").Where("student_id = ?", student.StudentId).Find(&orders).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	type OrderResponse struct {
		OrderID  uint `json:"order_id"`
		Quantity uint `json:"quantity"`
		Product  struct {
			ID       uint   `json:"product_id"`
			Name     string `json:"name"`
			Price    uint   `json:"price"`
			ImageURL string `json:"image_URL,omitempty"`
		} `json:"product"`
	}

	var response []OrderResponse
	for _, order := range orders {
		response = append(response, OrderResponse{
			OrderID:  order.OrderId,
			Quantity: order.Quantity,
			Product: struct {
				ID       uint   `json:"product_id"`
				Name     string `json:"name"`
				Price    uint   `json:"price"`
				ImageURL string `json:"image_URL,omitempty"`
			}{
				ID:       order.Product.ProductId,
				Name:     order.Product.Name,
				Price:    order.Product.Price,
				ImageURL: order.Product.ImageURL,
			},
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func ProductUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		Price       uint   `json:"price"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Quantity    uint   `json:"quantity"`
	}
	if request.Name == "" || request.Price <= 0 {
		log.Print("Missing request fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	idProduct, err := strconv.Atoi(mux.Vars(r)["id_product"])
	if err != nil {
		log.Print("conversion error")
		http.Error(w, "conversion errorr", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print("Invalid JSON", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()
	var product models.Product
	if err := database.DB.Where("product_id = ?", idProduct).First(&product).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	database.DB.Model(&product).Updates(map[string]interface{}{
		"Price":       request.Price,
		"Name":        request.Name,
		"Description": request.Description,
		"Quantity":    request.Quantity})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
func UpdateStatusProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Print("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	productId := mux.Vars(r)["id_product"]
	var product models.Product
	if err := database.DB.Where("product_id = ?", productId).First(&product).Error; err != nil {
		log.Print("Invalid credentials", err)
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusNotFound)
		return
	}
	isDeleted := false
	if !product.IsDeleted {
		isDeleted = true
	}
	database.DB.Model(&product).Updates(map[string]interface{}{
		"isDeleted": isDeleted})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
