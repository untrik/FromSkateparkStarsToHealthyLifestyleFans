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
	return *r
}
