package services

import (
	"testing"

	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/stretchr/testify/assert"
)

type mockUserModel struct {
	user *models.User
}

func TestUserService_Get(t *testing.T) {
	mockModelUser := newMockUserModel()
	mockUserServ := NewUserService(mockModelUser.user)
	user, err := mockUserServ.Get(1) // When connect to db, the newMockUserModel should be updated.
	if assert.Nil(t, err) && assert.NotNil(t, user) {
		assert.Equal(t, uint64(1), user.ID)
		assert.Equal(t, "", user.FirstName)
		assert.Equal(t, "", user.LastName)
	}
}

func newMockUserModel() *mockUserModel {
	return &mockUserModel{
		user: &models.User{
			Model:     models.Model{ID: 1},
			FirstName: "", // gofakeit.FirstName(),
			LastName:  "", // gofakeit.LastName(),
			Address:   "", // gofakeit.Address().Address,
			Email:     "", // gofakeit.Email(),
		},
	}
}
