package routes

import (
	"github.com/go-chi/chi"
	"github.com/mjande/user-auth-microservice/handlers"
)

func RegisterAuthRoutes(router chi.Router) {
	router.Post("/register", handlers.RegisterUser)
	router.Post("/login", handlers.LoginUser)
}
