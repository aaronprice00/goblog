package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aaronprice00/goblog/api/auth"
	"github.com/aaronprice00/goblog/api/response"
	"github.com/aaronprice00/goblog/api/util/formaterror"
	"github.com/aaronprice00/goblog/model"
	"golang.org/x/crypto/bcrypt"
)

// Login ... yup
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
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

	user.Prepare()
	if err = user.Validate("login"); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedErr := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusUnprocessableEntity, formattedErr)
		return
	}
	response.JSON(w, http.StatusOK, token)
}

// SignIn creates the token
func (server *Server) SignIn(email string, password string) (string, error) {
	var err error
	u := model.User{}
	user, err := u.ReadUserByEmail(server.DB, email)
	if err != nil {
		return "", err
	}
	err = model.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
