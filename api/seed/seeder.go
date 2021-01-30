package seed

import (
	"fmt"
	"log"

	"github.com/aaronprice00/goblog-mvc/api/model"
	"gorm.io/gorm"
)

var users = []model.User{
	{
		Username: "aaronprice00",
		Email:    "aaronprice00@gmail.com",
		Password: "pass123",
	},
	{
		Username: "phlesh",
		Email:    "phlesh@gmail.com",
		Password: "pass123",
	},
}

var posts = []model.Post{
	{
		Title:   "Title 1",
		Content: "Content 3",
	},
	{
		Title:   "Title 2",
		Content: "Content 4",
	},
}

// Load Drops, Migrates, Creates
func Load(db *gorm.DB) {

	var err error
	err = db.Migrator().DropTable(&model.Post{}, &model.User{})
	if err != nil {
		log.Fatalf("Could not drop table: %v", err)
	} else {
		fmt.Println("Dropped Tables")
	}

	err = db.AutoMigrate(&model.Post{}, &model.User{})
	if err != nil {
		log.Fatalf("Could not migrate table: %v", err)
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Fatalf("Could not seed User table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		if err := db.Create(&posts[i]).Error; err != nil {
			log.Fatalf("could not seed Post table: %v", err)
		}
	}
}
