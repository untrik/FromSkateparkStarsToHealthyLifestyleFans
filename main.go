package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/untrik/FromSkateToZOH/database"
	"github.com/untrik/FromSkateToZOH/middleware"
	"github.com/untrik/FromSkateToZOH/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	middleware.InitSecretKey()
	database.InitDB()
	r := routes.SetupRouter()
	http.ListenAndServe(":8081", &r)
}
