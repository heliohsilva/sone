package controllers

import (
	"api/src/auth"
	"api/src/db"
	"api/src/models"
	"api/src/repository"
	"api/src/responses"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	var post models.Post

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	post.UserID = uint(userID)

	if err = post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	database, err := db.Conn()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	defer database.Close()

	postRepository := repository.NewPostRepository(database)

	post.ID, err = postRepository.Create(post)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)

}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	database, err := db.Conn()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer database.Close()

	postRepository := repository.NewPostRepository(database)

	posts, err := postRepository.GetPosts(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	postID, err := strconv.ParseInt(parameters["postID"], 10, 64)

	fmt.Print(postID)

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

	postRepository := repository.NewPostRepository(database)

	post, err := postRepository.FetchByID(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)

	postID, err := strconv.ParseInt(parameters["postID"], 10, 64)

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

	postRepository := repository.NewPostRepository(database)

	postInDB, err := postRepository.FetchByID(postID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if postInDB.UserID != uint(userID) {
		responses.Error(w, http.StatusForbidden, errors.New("It is not possible to update a post written by other user"))
		return
	}

	var post models.Post

	if err = json.NewDecoder(r.Body).Decode(&post); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = postRepository.UpdatePost(postID, post); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)

	postID, err := strconv.ParseInt(parameters["postID"], 10, 64)

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

	postRepository := repository.NewPostRepository(database)

	postInDB, err := postRepository.FetchByID(postID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if postInDB.UserID != uint(userID) {
		responses.Error(w, http.StatusForbidden, errors.New("It is not possible to delete a post written by other user"))
		return
	}

	if err = postRepository.DeletePost(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
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

	postRepository := repository.NewPostRepository(database)

	posts, err := postRepository.FetchByUser(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func Like(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	postID, err := strconv.ParseInt(parameters["postID"], 10, 64)

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

	postRepository := repository.NewPostRepository(database)

	if err := postRepository.Like(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func Unlike(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	postID, err := strconv.ParseInt(parameters["postID"], 10, 64)

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

	postRepository := repository.NewPostRepository(database)

	if err := postRepository.Unlike(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
