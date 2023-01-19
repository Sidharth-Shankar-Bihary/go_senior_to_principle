package models

import (
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
	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
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
	err = db.Where("username=?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return
}
