package controller

import (
	"net/http"

	"github.com/aaronprice00/goblog-mvc/api/response"
)

// Home is the index handler
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "Welcome to the API")
}
