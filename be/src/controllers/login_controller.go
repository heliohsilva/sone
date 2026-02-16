package controllers

import (
	"api/src/auth"
	"api/src/db"
	"api/src/models"
	"api/src/repository"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	database, err := db.Conn()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer database.Close()

	userRepository := repository.NewUserRepository(database)
	userFetch, err := userRepository.GetByEmail(user.Email)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	err = security.Compare(userFetch.Password, user.Password)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
	}

	token, err := auth.GenerateToken(int64(userFetch.ID))

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, 200, token)

}
