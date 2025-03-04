package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/mjande/user-auth-microservice/database"
	"github.com/mjande/user-auth-microservice/models"
	"github.com/mjande/user-auth-microservice/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error while hashing password")
		return
	}

	query := "INSERT INTO users (email, password) VALUES ($1, $2);"
	_, err = database.DB.Exec(r.Context(), query, user.Email, hashedPassword)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error inserting user")
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

	query := "SELECT id, email, password FROM users WHERE email = $1;"
	row := database.DB.QueryRow(r.Context(), query, user.Email)

	var hashedPassword string
	err := row.Scan(&user.ID, &user.Email, &hashedPassword)
	if err != nil && err == pgx.ErrNoRows {
		utils.SendErrorResponse(w, http.StatusBadRequest, "This user does not exist")
		return
	} else if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}

	if !utils.CheckPassword(user.Password, hashedPassword) {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	var data map[string]string
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	data = map[string]string{
		"token":   token,
		"message": "Successfully logged in",
	}

	json.NewEncoder(w).Encode(data)
}
