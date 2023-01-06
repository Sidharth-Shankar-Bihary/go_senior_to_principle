package services

import (
	"errors"
	"html"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/auth"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(repo UserRepo, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (u *userService) CreateUser(req RegisterRequest) (resp *RegisterResponse, err error) {
	resp = &RegisterResponse{}
	user := models.User{}

	// valid request
	if err = validator.New().Struct(req); err != nil {
		resp.Err = err
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	user.Username = html.EscapeString(strings.TrimSpace(req.Username))
	user.Email = html.EscapeString(strings.TrimSpace(req.Email))
	if err = u.repo.CreateUser(&user); err != nil {
		resp.Err = err
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	resp.Status = http.StatusOK
	resp.Err = nil
	return
}

func (u *userService) VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *userService) GetUserToken(rds *redis.Client, req LoginRequest) (resp *LoginResponse, err error) {
	resp = &LoginResponse{}
	user := &models.User{}

	// valid request
	if err = validator.New().Struct(req); err != nil {
		resp.Err = err
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	user, err = u.repo.GetUserByName(req.Username)
	if err != nil || user.ID == 0 {
		err = errors.New("user does not exit")
		resp.Err = errors.New("user does not exit")
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	err = u.VerifyPassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		resp.Token = ""
		resp.Err = err
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	td, err := auth.GenerateToken(rds, user.ID)
	if err != nil {
		resp.Token = ""
		resp.Err = err
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	resp.Token = td.AccessToken
	resp.RefreshToken = td.RefreshToken
	resp.Status = http.StatusOK
	resp.Err = nil
	return
}

func (u *userService) GetUserByID(req GetUserRequest) (resp *GetUserResponse, err error) {
	resp = &GetUserResponse{}
	if req.ID <= 0 {
		resp.Status = http.StatusBadRequest
		resp.Err = "request ID is wrong"
		return
	}

	// valid request
	if err = validator.New().Struct(req); err != nil {
		resp.Err = "request something wrong"
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	user, err := u.repo.GetUserByID(uint(req.ID))
	if err != nil {
		resp.Err = err.Error()
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	resp.User = user
	resp.Status = http.StatusOK
	resp.Err = ""

	return
}
