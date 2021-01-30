package controllertest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aaronprice00/goblog-mvc/api/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatalf("Could not Refresh User Table, Error: %v \n", err)
	}
	samples := []struct {
		testID       int
		inputJSON    string
		statusCode   int
		username     string
		email        string
		errorMessage string
	}{
		{
			// sucessful
			testID:       1,
			inputJSON:    `{"username": "jcousteau", "email": "jcousteau@gmail.com", "password": "pass123"}`,
			statusCode:   201,
			username:     "jcousteau",
			email:        "jcousteau@gmail.com",
			errorMessage: "",
		},
		{
			testID:       2,
			inputJSON:    `{"username": "abuhlmann", "email": "jcousteau@gmail.com", "password": "pass123"}`,
			statusCode:   500,
			errorMessage: "Email Already Used",
		},
		{
			testID:       3,
			inputJSON:    `{"username": "jcousteau", "email": "abuhlmann@gmail.com", "password": "pass123"}`,
			statusCode:   500,
			errorMessage: "Username Already Taken",
		},
		{
			testID:       4,
			inputJSON:    `{"username": "abuhlman", "email": "abuhlmanngmail.com", "password": "pass123"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			testID:       5,
			inputJSON:    `{"username": "", "email": "abuhlmann@gmail.com", "password": "pass123"}`,
			statusCode:   422,
			errorMessage: "Required: Username",
		},
		{
			testID:       6,
			inputJSON:    `{"username": "abuhlmann", "email": "", "password": "pass123"}`,
			statusCode:   422,
			errorMessage: "Required: Email",
		},
		{
			testID:       7,
			inputJSON:    `{"username": "abuhlmann", "email": "abuhlmann@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required: Password",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		// Process response
		responseMap := make(map[string]interface{})
		if err = json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
			t.Errorf("Could not convert to JSON, Error: %v \n", err)
		}

		// Created User, check responding fields
		assert.Equal(t, v.statusCode, rr.Code)
		if v.statusCode == 201 {
			assert.Equal(t, v.username, responseMap["username"])
			assert.Equal(t, v.email, responseMap["email"])
		}

		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, v.errorMessage, responseMap["error"])
		}
		fmt.Printf("%v Finished w/ code: %v\n", v.testID, rr.Code)
	}
}

func TestGetUsers(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatalf("Could not refresh User Table, Error: %v \n", err)
	}
	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Could not seed Users, Error: %v \n", err)
	}
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("Error: %v \n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)

	var usersReceived []model.User
	if err = json.Unmarshal([]byte(rr.Body.String()), &usersReceived); err != nil {
		log.Fatalf("Could not convert to JSON, Error: %v \n", err)
	}
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, len(users), len(usersReceived))
}

func TestGetUserByID(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatalf("Could not refresh User table, Error: %v \n", err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Could not seed User, Error: %v \n", err)
	}

	userSample := []struct {
		id         string
		statusCode int
		username   string
		email      string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			username:   user.Username,
			email:      user.Email,
		}, {
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range userSample {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		if err = json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
			log.Fatalf("Could not convert to JSON, Error: %v \n", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Username, responseMap["username"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	var authEmail, authPassword string
	var authID uint

	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatalf("Could not refresh user table, Error: %v \n", err)
	}

	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Could not seed Users, Error: %v \n", err)
	}

	// Get only first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		authID = user.ID
		authEmail = user.Email
		authPassword = "pass123" // Unhashed
	}

	// Login the user and get authID
	token, err := server.SignIn(authEmail, authPassword)
	if err != nil {
		log.Fatalf("Could not login user, Error: %v \n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		testID         int
		id             string
		updateJSON     string
		statusCode     int
		updateUsername string
		updateEmail    string
		tokenGiven     string
		errorMessage   string
	}{
		{
			testID:         1,
			id:             strconv.Itoa(int(authID)),
			updateJSON:     `{"username": "willywonka", "email": "willy@wonkamail.com", "password": "pass123"}`,
			statusCode:     200,
			updateUsername: "willywonka",
			updateEmail:    "willy@wonkamail.com",
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			testID:       2,
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"username": "ebaker", "email": "erik@baker.com", "password": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required: Password",
		},
		{
			testID:       3,
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"username": "ebaker", "email": "erik@baker.com", "password": "pass123"}`,
			statusCode:   401,
			tokenGiven:   "", // Blank token
			errorMessage: "Unauthorized",
		},
		{
			testID:       4,
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"username": "ebaker", "email": "erik@baker.com", "password": "pass123"}`,
			statusCode:   401,
			tokenGiven:   "badtoken", // Bad token
			errorMessage: "Unauthorized",
		},
		{
			// email: albert@buhlmann.com belongs to user 2, should fail to update
			testID:       5,
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"username": "ebaker", "email": "albert@buhlmann.com", "password": "pass123"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Email Already Used",
		},
		{
			// username: abuhlmann belogs to user 2, should fail to update
			testID:       6,
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"username": "abuhlmann", "email": "erik@baker.com", "password": "pass123"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Username Already Taken",
		},
		{
			testID:       7,
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"username": "ebaker", "email": "erikbaker.com", "password": "pass123"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Invalid Email",
		},
		{
			testID:       8,
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"username": "", "email": "erik@baker.com", "password": "pass123"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required: Username",
		},
		{
			testID:       9,
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"username": "ebaker", "email": "", "password": "pass123"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required: Email",
		},
		{
			testID:     10,
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When user 2 is using user 1 token
			testID:       11,
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"username": "ebaker", "email": "erik@baker.com", "password": "pass123"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {
		var err error
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		if err := json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
			t.Errorf("Could not convert to JSON, Error: %v \n", err)
		}
		assert.Equal(t, v.statusCode, rr.Code)
		if v.statusCode == 200 {
			assert.Equal(t, v.updateUsername, responseMap["username"])
			assert.Equal(t, v.updateEmail, responseMap["email"])
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, v.errorMessage, responseMap["error"])
		}
		fmt.Printf("%v Finished w/ code: %v\n", v.testID, rr.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	var authEmail, authPassword string
	var authID uint

	var err error
	if err := refreshUserTable(); err != nil {
		log.Fatalf("Could not refresh User Table, Error: %v \n", err)
	}

	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Could not seed Users, Error: %v \n", err)
	}

	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		authID = user.ID
		authEmail = user.Email
		authPassword = "pass123" // unhashed password
	}
	// Login and grab auth token
	token, err := server.SignIn(authEmail, authPassword)
	if err != nil {
		log.Fatalf("Could not login, Error: %v \n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		testID       int
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			testID:       1,
			id:           strconv.Itoa(int(authID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			testID:       2,
			id:           strconv.Itoa(int(authID)),
			tokenGiven:   "", // Blank token
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			testID:       3,
			id:           strconv.Itoa(int(authID)),
			tokenGiven:   "badtoken", // bad token
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			testID:     4,
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400, // Bad Request
		},
		{
			// User 2 trying to use User 1 token
			testID:       5,
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range userSample {
		var err error
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage == "" {
			responseMap := make(map[string]interface{})
			if err = json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
				t.Errorf("Could not convert to JSON, Error: %v \n", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
		fmt.Printf("%v Finished w/ code: %v\n", v.testID, rr.Code)
	}
}
