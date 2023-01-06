package repos

import (
	"testing"

	"github.com/ramseyjiang/go_senior_to_principle/internal/api/services"
)

type mockUserRepo struct {
	// UserRepo User
}

func TestGetUserRepo(t *testing.T) {
	// req := GetUserRequest{ID: 1}
	// realRepo := NewUserRepo()
	// realResp, _ := realRepo.GetUser(req)
	//
	// mockRepoUser := newMockUserRepo()
	// mockResp, _ := mockRepoUser.GetUser(req)
	//
	// assert.Equal(t, mockResp, realResp)
	// assert.Equal(t, mockResp.User.ID, realResp.User.ID)
	// assert.Equal(t, mockResp.Status, realResp.Status)
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{}
}

// Mock Get function that replaces real userInterface
func (m *mockUserRepo) GetUser(req services.GetUserRequest) (resp *services.GetUserResponse, err error) {

	return resp, nil
}
