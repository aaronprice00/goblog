package controllertest

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aaronprice00/goblog-mvc/api/controller"
	"github.com/aaronprice00/goblog-mvc/api/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var server = controller.Server{}
var userInstance = model.User{}
var postInstance = model.Post{}

func TestMain(m *testing.M) {
	if err := godotenv.Load(os.ExpandEnv("../../.env")); err != nil {
		log.Fatalf("Could not get .env data %v \n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("T_DB_HOST"), os.Getenv("T_DB_PORT"), os.Getenv("T_DB_USER"), os.Getenv("T_DB_NAME"), os.Getenv("T_DB_PASSWORD"))
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Could not connect to database")
		log.Fatalf("Error: %v \n", err)
	} else {
		fmt.Println("Db connected :)")
	}
}

func refreshUserTable() error {
	var err error
	if err = server.DB.Migrator().DropTable(&model.User{}); err != nil {
		return err
	}
	if err = server.DB.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	log.Printf("Refreshed User table successfully")
	return nil
}

func seedOneUser() (model.User, error) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatal(err)
	}

	user := model.User{
		Username: "willywonka",
		Email:    "willy@wonkamail.com",
		Password: "pass123",
	}

	if err = server.DB.Model(&model.User{}).Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func seedUsers() ([]model.User, error) {
	users := []model.User{
		{
			Username: "jcousteau",
			Email:    "jaques@cousteau.com",
			Password: "pass123",
		},
		{
			Username: "abuhlmann",
			Email:    "albert@buhlmann.com",
			Password: "pass123",
		},
	}
	var err error
	if err = server.DB.Create(&users).Error; err != nil {
		return []model.User{}, err
	}
	return users, nil
}

func refreshUserAndPostTable() error {
	var err error
	if err = server.DB.Migrator().DropTable(&model.User{}, &model.Post{}); err != nil {
		return err
	}
	if err = server.DB.AutoMigrate(&model.User{}, &model.Post{}); err != nil {
		return err
	}
	log.Printf("Refreshed tables successfully")
	return nil
}

func seedOneUserAndOnePost() (model.Post, error) {
	var err error
	if err = refreshUserAndPostTable(); err != nil {
		return model.Post{}, err
	}
	user := model.User{
		Username: "jcousteau",
		Email:    "jaques@cousteau.com",
		Password: "pass123",
	}
	if err = server.DB.Create(&user).Error; err != nil {
		return model.Post{}, err
	}
	post := model.Post{
		Title:    "Under the sea",
		Content:  "Life is better, down where it's wetter",
		AuthorID: user.ID,
	}
	if err = server.DB.Create(&post).Error; err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func seedUsersAndPosts() ([]model.User, []model.Post, error) {

	var users = []model.User{
		{
			Username: "jcousteau",
			Email:    "jaques@cousteau.com",
			Password: "pass123",
		},
		{
			Username: "abuhlmann",
			Email:    "albert@buhlmann.com",
			Password: "pass123",
		},
	}
	var posts = []model.Post{
		{
			Title:   "We got no troubles",
			Content: "Life is the bubbles",
		},
		{
			Title:   "Compartment Schompartment",
			Content: "why not 16?",
		},
	}

	for i := range users {
		var err error
		if err = server.DB.Create(&users[i]).Error; err != nil {
			log.Fatalf("Could not seed users table: %v \n", err)
		}
		posts[i].AuthorID = users[i].ID

		if err = server.DB.Create(&posts[i]).Error; err != nil {
			log.Fatalf("Could not seed posts table: %v \n", err)
		}
	}
	return users, posts, nil
}
