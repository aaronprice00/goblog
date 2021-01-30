package api

import (
	"fmt"
	"log"
	"os"

	"github.com/aaronprice00/goblog-mvc/api/controller"
	"github.com/aaronprice00/goblog-mvc/api/seed"
	"github.com/joho/godotenv"
)

var server = controller.Server{}

// Run the REST server
func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load .env file %v", err)
	}

	server.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	seed.Load(server.DB)
	server.Run(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")))
}
