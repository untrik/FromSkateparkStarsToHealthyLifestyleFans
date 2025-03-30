package main

import (
	"net/http"

	"github.com/untrik/FromSkateToZOH/database"
	"github.com/untrik/FromSkateToZOH/routes"
)

func main() {
	database.InitDB()
	r := routes.SetupRouter()
	http.ListenAndServe(":8081", &r)
}
