package controllertest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatalf("Could not refresh user table, Error: %v \n", err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Could not seed user, Error: %v \n", err)
	}

	samples := []struct {
		testID       int
		email        string
		password     string
		errorMessage string
	}{
		{
			testID:       1,
			email:        user.Email,
			password:     "pass123", // must be non-hashed value
			errorMessage: "",
		},
		{
			testID:       2,
			email:        user.Email,
			password:     "wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			testID:       3,
			email:        "Wrong email",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {
		token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, errors.New(v.errorMessage), err)
		} else {
			assert.NotEqual(t, "", token)
		}
		fmt.Printf("%v Finished \n", v.testID)
	}
}

func TestLogin(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatalf("Could not refresh user table, Error: %v \n", err)
	}
	if _, err = seedOneUser(); err != nil {
		log.Fatalf("Could not seed user, Error: %v \n", err)
	}

	samples := []struct {
		testID       int
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			testID:       1,
			inputJSON:    `{"email": "willy@wonkamail.com", "password": "pass123"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			testID:       2,
			inputJSON:    `{"email": "jcousteau@gmail.com", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: "Incorrect Details",
		},
		{
			testID:       3,
			inputJSON:    `{"email": "banana@gmail.com", "password": "pass123"}`,
			statusCode:   422,
			errorMessage: "Incorrect Details",
		},
		{
			testID:       4,
			inputJSON:    `{"email": "jcousteaugmail.com", "password": "pass123"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			testID:       5,
			inputJSON:    `{"email": "", "password": "pass123"}`,
			statusCode:   422,
			errorMessage: "Required: Email",
		},
		{
			testID:       6,
			inputJSON:    `{"email": "jcousteau@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required: Password",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		// Log in successful, no data returned
		assert.Equal(t, v.statusCode, rr.Code)
		if v.statusCode != 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}

		// Login not successful, process response
		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			if err = json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
				t.Errorf("Could not convert to JSON, Error: %v \n", err)
			}
			assert.Equal(t, v.errorMessage, responseMap["error"])
		}
		fmt.Printf("%v Finished \n", v.testID)
	}
}
