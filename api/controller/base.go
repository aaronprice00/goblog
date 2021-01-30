package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aaronprice00/goblog-mvc/api/model"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Server holds our db and router objects
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize intitializes Server object with open db connection and routed Router
func (server *Server) Initialize(DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Cannot connect to database")
		log.Fatalln("Db Error: ", err)
	} else {
		fmt.Println("Db Connected")
	}

	server.DB.AutoMigrate(&model.User{}, &model.Post{})

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

// Run starts http Listen and Serve
func (server *Server) Run(addr string) {
	fmt.Println("Listening on port", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
