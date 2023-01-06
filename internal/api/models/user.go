package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User Password field "-" means it ensures the user’s password is not returned to the JSON response.
type User struct {
	Model
	Username string `gorm:"size:255;not null;" json:"username"`
	Password string `gorm:"size:255;not null;" json:"-"`
	Address  string `gorm:"size:1000" json:"address,omitempty"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
}

func (u *User) CreateUser(user *User, db *gorm.DB) (err error) {
	user.Password, err = u.HashPassword(user.Password)
	if err != nil {
		return err
	}

	if err = db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) HashPassword(pwd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (u *User) GetUserByID(userID uint64, db *gorm.DB) (user *User, err error) {
	user = &User{}
	user.ID = userID
	if result := db.Find(&user); result.Error != nil {
		return nil, result.Error
	}

	return
}

func (u *User) GetUserByName(username string, db *gorm.DB) (user *User, err error) {
	user = &User{}
	// Use where condition to get exactly result, if here use find only, using the login to check the current user will have wrong userInfo.
	if err = db.Where("username=?", username).Find(&user).Error; err != nil {
		return nil, err
	}
	return
}
