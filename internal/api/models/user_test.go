package models

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func TestCreateUser(t *testing.T) {
	mockUser := &User{}
	mockUser.Username = gofakeit.Username()
	mockUser.Email = gofakeit.Email()
	mockUser.Address = "test"
	mockUser.Password = "12345678"

	// wanted := mockUser.CreateUser(mockUser, db)
	// expected := ""
	//
	// assert.Equal(t, wanted, expected)
}
