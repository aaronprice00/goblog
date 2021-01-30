package controller

import (
	m "github.com/aaronprice00/goblog/api/middleware"
)

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", m.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", m.SetMiddlewareJSON(s.Login)).Methods("POST")

	// User Routes
	s.Router.HandleFunc("/users", m.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", m.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", m.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", m.SetMiddlewareJSON(m.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", m.SetMiddlewareJSON(s.DeleteUser)).Methods("DELETE")

	// Post Routes
	s.Router.HandleFunc("/posts", m.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	s.Router.HandleFunc("/posts", m.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", m.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", m.SetMiddlewareJSON(m.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", m.SetMiddlewareJSON(s.DeletePost)).Methods("DELETE")
}
