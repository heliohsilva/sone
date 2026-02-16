package controllers

import (
	"api/src/auth"
	"api/src/db"
	"api/src/models"
	"api/src/repository"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.PrepareUser("Create")

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

	userID, err := userRepository.CreateUser(user)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	user.ID = uint(userID)
	responses.JSON(w, http.StatusCreated, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNickname := strings.ToLower(r.URL.Query().Get("user"))

	fmt.Print(nameOrNickname)
	database, err := db.Conn()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer database.Close()

	userRepository := repository.NewUserRepository(database)

	users, err := userRepository.GetUsers(nameOrNickname)

	responses.JSON(w, 200, users)

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	database, err := db.Conn()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer database.Close()

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	userRepository := repository.NewUserRepository(database)

	user, err := userRepository.GetUser(userID)

	responses.JSON(w, 200, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	database, err := db.Conn()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}
	defer database.Close()

	vars := mux.Vars(r)

	userID, err := strconv.ParseInt(vars["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userTokenID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userTokenID != userID {
		responses.Error(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	var user models.User

	err = json.NewDecoder(r.Body).Decode(&user)
	user.ID = uint(userID)

	if err = user.PrepareUser("Editing"); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userRepository := repository.NewUserRepository(database)
	updatedRow, err := userRepository.UpdateUser(user)

	responses.JSON(w, 200, updatedRow)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userIDstr := vars["userID"]

	userID, err := strconv.ParseInt(userIDstr, 10, 64)

	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userTokenID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
	}

	if userTokenID != userID {
		responses.Error(w, http.StatusForbidden, errors.New("forbidden"))
	}

	database, err := db.Conn()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	userRepository := repository.NewUserRepository(database)

	if _, err := userRepository.DeleteUser(int(userID)); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, 200, nil)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := auth.ExtractUserID(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	userID, err := strconv.ParseInt(params["userID"], 10, 64)

	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if followerID == userID {
		responses.Error(w, http.StatusForbidden, errors.New("Is it impossible to follow yourself"))
		return
	}

	database, err := db.Conn()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	defer database.Close()

	userRepository := repository.NewUserRepository(database)

	if err = userRepository.Follow(userID, followerID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := auth.ExtractUserID(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)

	userID, err := strconv.ParseInt(parameters["userID"], 10, 64)

	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if followerID == userID {
		responses.Error(w, http.StatusForbidden, errors.New("Is not possible to follow yourself"))
		return
	}

	database, err := db.Conn()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	defer database.Close()

	userRepository := repository.NewUserRepository(database)

	if err = userRepository.UnfollowUser(userID, followerID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}

func FetchFollowers(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userID, err := strconv.ParseInt(parameters["userID"], 10, 64)

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

	followers, err := userRepository.FetchFollowers(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func FetchFollowing(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userID, err := strconv.ParseInt(parameters["userID"], 10, 64)

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

	users, err := userRepository.FetchFollowing(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userIDOnToken, err := auth.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)
	userID, err := strconv.ParseInt(parameters["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if userIDOnToken != userID {
		responses.Error(w, http.StatusForbidden, errors.New("Is it not possible to update another's user password"))
		return
	}

	var password models.Password

	err = json.NewDecoder(r.Body).Decode(&password)

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
	currentPassword, err := userRepository.FetchPassword(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.Compare(currentPassword, password.Old); err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Wrong password"))
		return
	}

	hashedPassword, err := security.Hash(password.New)

	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = userRepository.UpdatePassword(userID, string(hashedPassword)); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
