package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/aaronprice00/goblog/api/auth"
	"github.com/aaronprice00/goblog/api/response"
	"github.com/aaronprice00/goblog/api/util/formaterror"
	"github.com/aaronprice00/goblog/model"
	"github.com/gorilla/mux"
)

// CreateUser escapes, validates, authenticates before asking model to create
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := model.User{}
	if err = json.Unmarshal(body, &user); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Trim whitespaces and escape Username and Email
	user.Prepare()

	if err = user.Validate(""); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.CreateUser(server.DB)
	if err != nil {
		formattedErr := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formattedErr)
		return
	}

	// Location response header indicates the URL to redirect a page to
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	response.JSON(w, http.StatusCreated, userCreated)

}

// GetUsers asks model for list then responds with JSON
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	users, err := user.ReadAllUsers(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, users)
}

// GetUser grabs ID from the URL before asking model for the User
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := model.User{}
	userReceived, err := user.ReadUserByID(server.DB, uint(uid))
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	response.JSON(w, http.StatusOK, userReceived)
}

// UpdateUser grabs id from url, escapes, validates, and authenticates before asking model to update
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := model.User{}
	if err = json.Unmarshal(body, &user); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint(uid) {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	if err = user.Validate("update"); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateUser(server.DB, uint(uid))
	if err != nil {
		formattedErr := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formattedErr)
		return
	}
	response.JSON(w, http.StatusOK, updatedUser)
}

// DeleteUser pulls id from url and authenticates before asking model to delete responds via http JSON
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := model.User{}
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	if tokenID != 0 && tokenID != uint(uid) {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}

	if _, err := user.DeleteUser(server.DB, uint(uid)); err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	response.JSON(w, http.StatusNoContent, "")
}
