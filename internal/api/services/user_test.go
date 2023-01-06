package services

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
)

type mockUserRepo struct {
	user *models.User
}

func TestGetUserByID(t *testing.T) {
	// mockModelUser := newMockUserRepo()
	// req := &GetUserRequest{
	// 	ID: 1,
	// }
	//
	// mockUserServ := NewUserService()
	// user, err := mockUserServ.GetUserByID(*req) // When connect to db, the newMockUserModel should be updated.
	// if assert.Nil(t, err) && assert.NotNil(t, user) {
	// 	assert.Equal(t, uint64(1), mockModelUser.ID)
	// }
}

func newMockUserRepo() *models.User {
	return &models.User{
		Model:    models.Model{ID: 1},
		Username: gofakeit.FirstName(),
		Address:  gofakeit.Address().Address,
		Email:    gofakeit.Email(),
	}
}
