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

// CreatePost verifies validates and authorizes before creating it
func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := model.Post{}
	if err = json.Unmarshal(body, &post); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.Prepare()
	if err = post.Validate(); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != post.AuthorID {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	postCreated, err := post.CreatePost(server.DB)
	if err != nil {
		formattedErr := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formattedErr)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	response.JSON(w, http.StatusCreated, postCreated)
}

// GetPosts pulls post from model and responds via JSON
func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	post := model.Post{}
	posts, err := post.ReadAllPosts(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, posts)
}

// GetPost pulls id from the URL and asks model for post
func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post := model.Post{}

	postReceived, err := post.ReadPostByID(server.DB, uint(pid))
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, postReceived)
}

// UpdatePost pulls id from url escapes, validates, and authenticates before asking model to update
func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Is Post ID Valid?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is auth token valid? get user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Does Post Exist?
	p := model.Post{}
	post, err := p.ReadPostByID(server.DB, uint(pid))
	if err != nil {
		response.ERROR(w, http.StatusNotFound, err)
		return
	}

	// Don't allow user to update another user's post
	if uid != post.AuthorID {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Read the POST body data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing data
	postUpdate := model.Post{}
	if err = json.Unmarshal(body, &postUpdate); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Verify Post user is the same as Token user
	if uid != postUpdate.AuthorID {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	postUpdate.Prepare()

	if err = postUpdate.Validate(); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.ID = post.ID // Important to ensure the model knows which post row to update

	postUpdated, err := postUpdate.UpdatePost(server.DB)
	if err != nil {
		formattedErr := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formattedErr)
		return
	}

	response.JSON(w, http.StatusOK, postUpdated)
}

// DeletePost pulls id from URL, authenticates and asks model to delete
func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Is supplied post ID a number?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Get Post
	p := model.Post{}
	post, err := p.ReadPostByID(server.DB, uint(pid))
	if err != nil {
		response.ERROR(w, http.StatusNotFound, err)
		return
	}

	//Does this Post belong to this user?
	if uid != post.AuthorID {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Do the Delete
	if _, err := post.DeletePost(server.DB, uint(pid)); err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	response.JSON(w, http.StatusNoContent, "")
}
