package routes

import (
	"github.com/gorilla/mux"
	"github.com/untrik/FromSkateToZOH/controllers"
)

func SetupRouter() mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	return *r
}
