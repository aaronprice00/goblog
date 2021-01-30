package middleware

import (
	"errors"
	"net/http"

	"github.com/aaronprice00/goblog-mvc/api/auth"
	"github.com/aaronprice00/goblog-mvc/api/response"
)

// SetMiddlewareJSON sets content type of the response to JSON
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication checks to see if token is valid, if not set appropriate response
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
