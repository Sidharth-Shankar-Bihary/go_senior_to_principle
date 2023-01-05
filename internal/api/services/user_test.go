package services

import (
	"testing"

	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/stretchr/testify/assert"
)

type mockUserModel struct {
	user *models.User
}

func TestUserServiceGet(t *testing.T) {
	mockModelUser := newMockUserModel()
	mockUserServ := NewUserService(mockModelUser.user)
	user, err := mockUserServ.Get(1) // When connect to db, the newMockUserModel should be updated.
	if assert.Nil(t, err) && assert.NotNil(t, user) {
		assert.Equal(t, uint64(1), user.ID)
		assert.Equal(t, "", user.Username)
	}
}

func newMockUserModel() *mockUserModel {
	return &mockUserModel{
		user: &models.User{
			Model:    models.Model{ID: 1},
			Username: "", // gofakeit.FirstName(),
			Address:  "", // gofakeit.Address().Address,
			Email:    "", // gofakeit.Email(),
		},
	}
}
