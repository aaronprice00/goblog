package modeltest

import (
	"log"
	"testing"

	"github.com/aaronprice00/goblog-mvc/api/model"
	"github.com/stretchr/testify/assert"
)

func TestFindAllPosts(t *testing.T) {
	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh User and Post tables Error: %v \n", err)
	}
	if _, _, err = seedUsersAndPosts(); err != nil {
		log.Fatalf("Could not seed users and posts Error: %v \n", err)
	}
	posts, err := postInstance.ReadAllPosts(server.DB)
	if err != nil {
		t.Errorf("Could not find posts Error: %v \n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)
}

func TestCreatePost(t *testing.T) {
	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh user and post tables Error: %v \n", err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Could not seed user Error: %v \n", err)
	}

	newPost := model.Post{
		Title:    "New Title",
		Content:  "New Content",
		AuthorID: user.ID,
	}
	// Cant set "ID: 1" in struct litteral above because of "gorm.Model" in our Post Struct
	newPost.ID = 1

	savedPost, err := newPost.CreatePost(server.DB)
	if err != nil {
		t.Errorf("Could not Save Post Error: %v \n", err)
		return
	}
	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.Title, savedPost.Title)
	assert.Equal(t, newPost.Content, savedPost.Content)
	assert.Equal(t, newPost.AuthorID, savedPost.AuthorID)
}

func TestGetPostByID(t *testing.T) {
	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh User and post table Error: %v \n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Could not seed user and post Error: %v \n", err)
	}
	foundPost, err := postInstance.ReadPostByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("Could not find post Error: %v \n", err)
		return
	}
	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.Title, post.Title)
	assert.Equal(t, foundPost.Content, post.Content)
}

func TestUpdatePost(t *testing.T) {
	var err error
	if err = refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh User and Post table Error: %v \n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Could not seed user and post Error: %v \n", err)
	}

	postUpdate := model.Post{
		Title:    "Title 1",
		Content:  "Content 1",
		AuthorID: post.AuthorID,
	}
	// Cant set "ID: 1" in struct litteral above because of "gorm.Model" in our Post Struct
	postUpdate.ID = 1

	updatedPost, err := postUpdate.UpdatePost(server.DB)
	if err != nil {
		t.Errorf("Could not update post Error: %v \n", err)
		return
	}

	assert.Equal(t, updatedPost.ID, postUpdate.ID)
	assert.Equal(t, updatedPost.Title, postUpdate.Title)
	assert.Equal(t, updatedPost.Content, postUpdate.Content)
	assert.Equal(t, updatedPost.AuthorID, postUpdate.AuthorID)
}

func TestDeletePost(t *testing.T) {
	var err error
	if err := refreshUserAndPostTable(); err != nil {
		log.Fatalf("Could not refresh user and post table Error: %v \n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Could not seed user and post Error: %v \n", err)
	}
	isDeleted, err := postInstance.DeletePost(server.DB, post.ID)
	if err != nil {
		t.Errorf("Could not delete the post Error: %v \n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
