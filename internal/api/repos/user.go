package repos

import (
	"net/http"

	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/services"
)

type UserInterface interface {
	GetUserByID(req GetUserRequest) (*GetUserResponse, error)
}

type UserRepo struct {
	UserRepo UserInterface
}

type GetUserRequest struct {
	ID uint64
}

type GetUserResponse struct {
	User   *models.User `json:"user,omitempty"`
	Err    error        `json:"err,omitempty"`
	Status int          `json:"status,omitempty"`
}

// NewUserRepo creates a new NewUserRepo with the given user.
func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// GetUserByID just retrieves user using User Model, here can be additional logic for processing data retrieved by Models
func (u *UserRepo) GetUserByID(req GetUserRequest) (*GetUserResponse, error) {
	resp := new(GetUserResponse)          // The other way is resp := &GetUserResponse{}
	userServ := new(services.UserService) // The other way is userServ := &services.UserService{}

	// do some logic, then use userService to access data
	user, err := userServ.Get(req.ID)
	if err != nil {
		return nil, err
	}

	resp.User = user
	resp.Err = err
	resp.Status = http.StatusOK

	return resp, nil
}
