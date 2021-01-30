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

func TestCreatePost(t *testing.T) {
	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh User and Post tables, Error: %v \n", err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Could not seed User, Error: %v \n", err)
	}
	token, err := server.SignIn(user.Email, "pass123")
	if err != nil {
		log.Fatalf("Could not login, Error: %v \n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		testID       int
		inputJSON    string
		statusCode   int
		title        string
		content      string
		authorID     uint
		tokenGiven   string
		errorMessage string
	}{
		{
			// sucessful
			testID:       1,
			inputJSON:    `{"title": "Title 1", "content": "Content 1", "authorID": 1}`,
			statusCode:   201,
			title:        "Title 1",
			content:      "Content 1",
			authorID:     user.ID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			testID:       2,
			inputJSON:    `{"title": "Title 1", "content": "Content 2", "authorID": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Used",
		},
		{
			testID:       3,
			inputJSON:    `{"title": "Title 2", "content": "Content 2", "authorID": 1}`,
			statusCode:   401,
			tokenGiven:   "", // Blank Token
			errorMessage: "Unauthorized",
		},
		{
			testID:       4,
			inputJSON:    `{"title": "Title 2", "content": "Content 2", "authorID": 1}`,
			statusCode:   401,
			tokenGiven:   "badtoken", // Bad Token
			errorMessage: "Unauthorized",
		},
		{
			testID:       5,
			inputJSON:    `{"title": "", "content": "Content 2", "authorID": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required: Title",
		},
		{
			testID:       6,
			inputJSON:    `{"title": "Title 2", "content": "", "authorID": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required: Content",
		},
		{
			testID:       7,
			inputJSON:    `{"title": "Title 2", "content": "Content 2"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required: Author",
		},
		{
			// When user 2 uses user 1 token
			testID:       8,
			inputJSON:    `{"title": "Title 2", "content": "Content 2", "authorID": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreatePost)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		if err = json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
			t.Errorf("Could not convert to JSON, Error: %v \n", err)
		}

		assert.Equal(t, v.statusCode, rr.Code)
		if v.statusCode == 201 {
			assert.Equal(t, v.title, responseMap["title"])
			assert.Equal(t, v.content, responseMap["content"])
			assert.Equal(t, float64(v.authorID), responseMap["authorID"])
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, v.errorMessage, responseMap["error"])
		}
		fmt.Printf("%v Finished w/ code: %v\n", v.testID, rr.Code)
	}
}

func TestGetPosts(t *testing.T) {
	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh User and Post table, Error: %v \n", err)
	}
	_, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Could not seed users and posts, Error: %v \n", err)
	}

	req, err := http.NewRequest("GET", "/posts/", nil)
	if err != nil {
		t.Errorf("Error: %v \n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetPosts)
	handler.ServeHTTP(rr, req)

	var postsReceived []model.Post
	if err = json.Unmarshal([]byte(rr.Body.String()), &postsReceived); err != nil {
		t.Errorf("Could not convert to json, Error: %v \n", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, len(posts), len(postsReceived))
}

func TestGetPostByID(t *testing.T) {
	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not seed users and posts, Error: %v \n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Could not seed user and post, Error: %v \n", err)
	}

	samples := []struct {
		testID     int
		id         string
		statusCode int
		title      string
		content    string
		authorID   uint
	}{
		{
			testID:     1,
			id:         strconv.Itoa(int(post.ID)),
			statusCode: 200,
			title:      post.Title,
			content:    post.Content,
			authorID:   post.AuthorID,
		},
		{
			testID:     2,
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetPost)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		if err = json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
			t.Errorf("Could not convert to JSON, Error: %v \n", err)
		}
		assert.Equal(t, v.statusCode, rr.Code)

		if v.statusCode == 200 {
			assert.Equal(t, v.title, responseMap["title"])
			assert.Equal(t, v.content, responseMap["content"])
			assert.Equal(t, float64(v.authorID), responseMap["authorID"])
		}
		fmt.Printf("%v Finished w/ code: %v\n", v.testID, rr.Code)
	}
}

func TestUpdatePost(t *testing.T) {
	var postUserEmail, postUserPassword string
	var authPostAuthorID uint
	var authPostID uint

	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh user and post tables, Error: %v \n", err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Could not seed users and posts, Error: %v \n", err)
	}

	// grab first only the first user, wouldn't something like users[1].Email work?
	for _, u := range users {
		if u.ID == 2 {
			continue
		}
		postUserEmail = u.Email
		postUserPassword = "pass123" // unhashed
	}
	token, err := server.SignIn(postUserEmail, postUserPassword)
	if err != nil {
		t.Errorf("Could not login the user, Error: %v \n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// grab first only the first post, wouldn't something like posts[1].Email work?
	for _, p := range posts {
		if p.ID == 2 {
			continue
		}
		authPostID = p.ID
		authPostAuthorID = p.AuthorID
	}

	var samples = []struct {
		testID       int
		id           string
		updateJSON   string
		statusCode   int
		title        string
		content      string
		authorID     uint
		tokenGiven   string
		errorMessage string
	}{
		{
			// successful update
			testID:       1,
			id:           strconv.Itoa(int(authPostID)),
			updateJSON:   `{"title": "Title 1", "content": "Content 1", "authorID": 1}`,
			statusCode:   200,
			title:        "Title 1",
			content:      "Content 1",
			authorID:     authPostAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// blank token
			testID:       2,
			id:           strconv.Itoa(int(authPostID)),
			updateJSON:   `{"title": "Title 2", "content": "Content 2", "authorID": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// bad token
			testID:       3,
			id:           strconv.Itoa(int(authPostID)),
			updateJSON:   `{"title": "Title 2", "content": "Content 2", "authorID": 1}`,
			statusCode:   401,
			tokenGiven:   "badtoken",
			errorMessage: "Unauthorized",
		},
		{
			// no title
			testID:       4,
			id:           strconv.Itoa(int(authPostID)),
			updateJSON:   `{"title": "", "content": "Content 2", "authorID": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required: Title",
		},
		{
			// duplicate title, must be unique
			testID:       5,
			id:           strconv.Itoa(int(authPostID)),
			updateJSON:   `{"title": "Compartment Schompartment", "content": "Content 2", "authorID": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Used",
		},
		{
			// no content
			testID:       6,
			id:           strconv.Itoa(int(authPostID)),
			updateJSON:   `{"title": "Title 2", "content": "", "authorID": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required: Content",
		},
		{
			// no authorID
			testID:       7,
			id:           strconv.Itoa(int(authPostID)),
			updateJSON:   `{"title": "Title 2", "content": "Content 2"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			// Bad request
			testID:     8,
			id:         "unknown",
			statusCode: 400,
		},
		{
			// user's tokenID mismatch authorID
			testID:       9,
			id:           strconv.Itoa(int(authPostID)),
			updateJSON:   `{"title": "Title 2", "content": "Content 2", "authorID": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {
		var err error
		req, err := http.NewRequest("PUT", "/posts", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdatePost)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		if err = json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
			t.Errorf("Could not convert to JSON, Error: %v \n", err)
		}
		assert.Equal(t, v.statusCode, rr.Code)

		if rr.Code == 200 {
			assert.Equal(t, v.title, responseMap["title"])
			assert.Equal(t, v.content, responseMap["content"])
			assert.Equal(t, float64(v.authorID), responseMap["authorID"])
		}

		// What about the 400 error? Do we need to check that one?
		if rr.Code == 401 || rr.Code == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, v.errorMessage, responseMap["error"])
		}
		fmt.Printf("%v Finished w/ code: %v\n", v.testID, rr.Code)
	}
}

func TestDeletePost(t *testing.T) {
	var postUserEmail, postUserPassword string
	var authPostID, authPostAuthorID uint

	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh user and post tables, Error: %v \n", err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Could not seed user and post, Error: %v \n", err)
	}

	// Grab the first user
	for _, u := range users {
		if u.ID == 2 {
			continue
		}
		postUserEmail = u.Email
		postUserPassword = "pass123" // unhashed
	}

	// Grab the first post
	for _, p := range posts {
		if p.ID == 2 {
			continue
		}
		authPostID = p.ID
		authPostAuthorID = p.AuthorID
	}

	// login the user
	token, err := server.SignIn(postUserEmail, postUserPassword)
	if err != nil {
		t.Errorf("Could not login the user, Error: %v \n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	sample := []struct {
		testID       int
		id           string
		authorID     uint
		statusCode   int
		tokenGiven   string
		errorMessage string
	}{
		{
			// successful delete
			testID:       1,
			id:           strconv.Itoa(int(authPostID)),
			authorID:     authPostAuthorID,
			statusCode:   204,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// blank token
			testID:       2,
			id:           strconv.Itoa(int(authPostID)),
			authorID:     authPostAuthorID,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// bad token
			testID:       3,
			id:           strconv.Itoa(int(authPostID)),
			authorID:     authPostAuthorID,
			statusCode:   401,
			tokenGiven:   "badtoken",
			errorMessage: "Unauthorized",
		},
		{
			// bad request
			testID:     4,
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// ??
			testID:       5,
			id:           strconv.Itoa(int(1)),
			authorID:     1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range sample {
		var err error

		req, err := http.NewRequest("DELETE", "/posts", nil)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeletePost)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			if err = json.Unmarshal([]byte(rr.Body.String()), &responseMap); err != nil {
				t.Errorf("Could not convert to Json %v \n", err)
			}
			assert.Equal(t, v.errorMessage, responseMap["error"])
		}
		fmt.Printf("%v Finished w/ code: %v\n", v.testID, rr.Code)
	}
}
