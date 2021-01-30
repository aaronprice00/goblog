package modeltest

import (
	"log"
	"testing"

	"github.com/aaronprice00/goblog-mvc/api/model"
	"github.com/stretchr/testify/assert"
)

func TestFindAllUsers(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatalln(err)
	}
	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Could not seed Users, Error: %v \n", err)
	}

	usersReceived, err := userInstance.ReadAllUsers(server.DB)
	if err != nil {
		t.Errorf("Could not get users Error: %v \n", err)
		return
	}
	assert.Equal(t, len(users), len(*usersReceived))
}

func TestSaveUser(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatal(err)
	}
	newUser := model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "pass123",
	}
	// Cant set "ID: 1" in struct litteral above because of "gorm.Model" in our User Struct
	newUser.ID = 1

	savedUser, err := newUser.CreateUser(server.DB)
	if err != nil {
		t.Errorf("Could not get user Error: %v \n", err)
		return
	}
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Username, savedUser.Username)
	assert.Equal(t, newUser.Email, savedUser.Email)
}

func TestGetUserByID(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Could not seed User tabel %v \n", err)
		return
	}
	foundUser, err := userInstance.ReadUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Could not read user Error: %v \n", err)
		return
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Username, user.Username)
	assert.Equal(t, foundUser.Email, user.Email)
}

func TestUpdateUser(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Could not seed user: %v \n", err)
	}

	userUpdate := model.User{
		Username: "willywonka",
		Email:    "willy@wonkamail.com",
		Password: "pass123",
	}
	// Cant set "ID: 1" in struct litteral above because of "gorm.Model" in our User Struct
	userUpdate.ID = 1

	userUpdated, err := userUpdate.UpdateUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("Could not create user Error: %v \n", err)
	}

	assert.Equal(t, user.ID, userUpdated.ID)
	assert.Equal(t, user.Username, userUpdated.Username)
	assert.Equal(t, user.Email, userUpdated.Email)
}

func TestDeleteUser(t *testing.T) {
	var err error
	if err = refreshUserTable(); err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Could not seed user Error: %v \n", err)
	}

	isDeleted, err := userInstance.DeleteUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("Could not delete user Error: %v \n", err)
	}

	assert.Equal(t, isDeleted, int64(1))
}
