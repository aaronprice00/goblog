package model

import (
	"errors"
	"html"
	"log"
	"strings"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User holds our User; gorm.Model contains ID, CreatedAt, DeletedAt, and UpdatedA details
type User struct {
	gorm.Model
	Username string `gorm:"size:100;not null;unique;" json:"username"`
	Email    string `gorm:"size:100;not null;unique;" json:"email"`
	Password string `gorm:"size:100;not null;" json:"password"`
}

// Hash encrypts the supplied password returns hash and error
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword compares hashed and unhashed password returns nil on success, error on failure
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave hashes password and stores hashed to User object
func (u *User) BeforeSave(*gorm.DB) error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Prepare Escapes and trims posted Username & Email
func (u *User) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

// Validate checks required fields on specified actions
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Required: Username")
		}
		if u.Password == "" {
			return errors.New("Required: Password")
		}
		if u.Email == "" {
			return errors.New("Required: Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required: Password")
		}
		if u.Email == "" {
			return errors.New("Required: Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	default:
		if u.Username == "" {
			return errors.New("Required: Username")
		}
		if u.Password == "" {
			return errors.New("Required: Password")
		}
		if u.Email == "" {
			return errors.New("Required: Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

// CreateUser Inserts user into db returns User and error
func (u *User) CreateUser(db *gorm.DB) (*User, error) {
	if err := db.Create(&u).Error; err != nil {
		return &User{}, err
	}
	// Should do a fresh pull here
	return u, nil
}

// ReadAllUsers returns all results from User table
func (u *User) ReadAllUsers(db *gorm.DB) (*[]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return &[]User{}, err
	}
	return &users, nil
}

// ReadUserByID queries User table by ID and returns the matching user
func (u *User) ReadUserByID(db *gorm.DB, uid uint) (*User, error) {
	var err error
	if err = db.Take(&u, uid).Error; err != nil {
		return &User{}, err
	}
	// User not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// ReadUserByEmail queries User table by email returns the matching user
func (u *User) ReadUserByEmail(db *gorm.DB, email string) (*User, error) {
	var err error
	if err = db.Where("email = ?", email).Take(&u).Error; err != nil {
		return &User{}, err
	}
	// User not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// UpdateUser saves fields to User row at supplied ID
func (u *User) UpdateUser(db *gorm.DB, uid uint) (*User, error) {
	var err error
	// Hash the password
	if err = u.BeforeSave(db); err != nil {
		log.Fatalln(err)
	}
	res := db.Model(&User{}).Where("id = ?", uid).Updates(map[string]interface{}{
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
	})
	if err = res.Error; err != nil {
		return &User{}, err
	}

	// Grab a fresh copy
	if err = db.Take(&u, uid).Error; err != nil {
		return &User{}, err
	}

	return u, nil
}

// DeleteUser sets User row inactive (will need to purge), the return int is used for Testing suite to check isDeleted = 1
func (u *User) DeleteUser(db *gorm.DB, uid uint) (int64, error) {
	var err error
	res := db.Delete(&u, uid)
	if err = res.Error; err != nil {
		return res.RowsAffected, err
	}
	return res.RowsAffected, nil
}
