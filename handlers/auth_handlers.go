package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mjande/user-auth-microservice/database"
	"github.com/mjande/user-auth-microservice/models"
	"github.com/mjande/user-auth-microservice/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	query := "INSERT INTO users VALUES (?, ?)"
	_, err = database.DB.Exec(r.Context(), query, user.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := map[string]string{
		"message": "User created successfully",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	query := "SELECT id, email, password FROM users WHERE email = ?"
	row := database.DB.QueryRow(r.Context(), query, user.Email)

	var id int
	var hashedPassword string
	row.Scan(&id, &hashedPassword)

	if !utils.CheckPassword(user.Password, hashedPassword) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(id)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := map[string]string{
		"token": token,
	}

	json.NewEncoder(w).Encode(data)
}
