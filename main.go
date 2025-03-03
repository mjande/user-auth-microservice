package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mjande/user-auth-microservice/database"
	"github.com/mjande/user-auth-microservice/routes"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	err = database.InitDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer database.DB.Close()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("CLIENT_URL")},
		AllowedMethods: []string{"POST"},
	}))

	routes.RegisterAuthRoutes(router)

	log.Printf("Auth service listening on port %s", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if err != nil {
		log.Fatal(err)
	}
}
