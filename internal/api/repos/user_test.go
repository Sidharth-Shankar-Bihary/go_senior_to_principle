package repos

import (
	"net/http"
	"testing"

	"github.com/ramseyjiang/go_senior_to_principle/internal/api/services"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	UserRepo UserInterface
}

func TestUserRepo_GetUserByID(t *testing.T) {
	req := GetUserRequest{ID: 1}
	realRepo := NewUserRepo()
	realResp, _ := realRepo.GetUserByID(req)

	mockRepoUser := newMockUserRepo()
	mockResp, _ := mockRepoUser.GetUserByID(req)

	assert.Equal(t, mockResp, realResp)
	assert.Equal(t, mockResp.User.ID, realResp.User.ID)
	assert.Equal(t, mockResp.Status, realResp.Status)
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{}
}

// Mock Get function that replaces real userInterface
func (m *mockUserRepo) GetUserByID(req GetUserRequest) (*GetUserResponse, error) {
	resp := new(GetUserResponse)

	userServ := new(services.UserService)
	user, err := userServ.Get(req.ID)
	if err != nil {
		return nil, err
	}

	resp.User = user
	resp.Err = nil
	resp.Status = http.StatusOK

	return resp, nil
}
