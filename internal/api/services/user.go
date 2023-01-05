package services

import (
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
)

type User struct {
	user models.User
}

// NewUserService creates a NewUserService with the given user. The user_test will use this method also.
func NewUserService(user *models.User) *User {
	return &User{*user}
}

// Get just retrieves user using User Model, here can be additional logic for processing data retrieved by Models
func (u *User) Get(id uint64) (*models.User, error) {
	return u.user.Get(id) // No additional logic, just return the query result
}
