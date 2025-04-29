package routes

import (
	"github.com/gorilla/mux"
	"github.com/untrik/FromSkateToZOH/controllers"
	"github.com/untrik/FromSkateToZOH/middleware"
)

func SetupRouter() mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.CreateStudent).Methods("POST")
	r.HandleFunc("/admin/register", middleware.JWTMiddlewareAdmin(controllers.CreateAdmin)).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/event/add", middleware.JWTMiddlewareAdmin(controllers.CreateEvent)).Methods("POST")
	r.HandleFunc("/event/{id_event}/student", middleware.JWTMiddlewareStudent(controllers.RegistrationForTheEvent)).Methods("POST")
	r.HandleFunc("/reward/event/{id_event}/student/{id_student}", middleware.JWTMiddlewareAdmin(controllers.AddReward)).Methods("POST")
	r.HandleFunc("/event/{id_event}/status/student/{id_student}/change", middleware.JWTMiddlewareAdmin(controllers.ChangeStatusStudent)).Methods("PUT")
	r.HandleFunc("/all/events", controllers.GetAllActiveEvent).Methods("GET")
	r.HandleFunc("/event/{id_event}", middleware.JWTMiddlewareAdmin(controllers.DeleteEvent)).Methods("DELETE")
	r.HandleFunc("/all/product", controllers.GetAllProduct).Methods("GET")
	r.HandleFunc("/product/add", middleware.JWTMiddlewareAdmin(controllers.AddProduct)).Methods("POST")
	r.HandleFunc("/order/{id_product}/add", middleware.JWTMiddlewareStudent(controllers.AddOrder)).Methods("POST")
	r.HandleFunc("/all/orders", middleware.JWTMiddlewareStudent(controllers.GetAllOrders)).Methods("GET")
	r.HandleFunc("/event/{id_event}/all/participants", middleware.JWTMiddlewareAdmin(controllers.GetAllParticipants)).Methods("GET")
	r.HandleFunc("/product/{id_product}", middleware.JWTMiddlewareAdmin(controllers.ProductUpdate)).Methods("PUT")
	r.HandleFunc("/product/{id_product}/status/update", middleware.JWTMiddlewareAdmin(controllers.UpdateStatusProduct)).Methods("PUT")
	return *r
}
