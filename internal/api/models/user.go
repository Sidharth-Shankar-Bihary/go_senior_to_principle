package models

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Username string `gorm:"column:username,size:255;not null;" json:"username"`
	Password string `gorm:"size:255;not null;" json:"-"`
	Address  string `gorm:"column:address" json:"address"`
	Email    string `gorm:"column:email;size:255;not null;unique" json:"email"`
}

var user User

// NewUser creates a new User
func NewUser() *User {
	return &User{}
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))

	return nil
}

func (u *User) Get(id uint64) (*User, error) {
	user.ID = id
	// err := env.Config.DB.Where("id = ?", id). // Do the query
	// 	First(&user). // Make it scalar
	// 	Error // retrieve error or null
	// return &user, err
	return &user, nil
}
