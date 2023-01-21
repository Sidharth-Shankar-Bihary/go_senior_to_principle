package services

import (
	"github.com/go-redis/redis/v8"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/repos"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(req RegisterRequest) (*RegisterResponse, error)
	VerifyPassword(password, hashedPassword string) error
	GetUserToken(rds *redis.Client, req LoginRequest) (*LoginResponse, error)
	GetUserByID(req GetUserRequest) (*GetUserResponse, error)
}

type userService struct {
	repo   repos.UserRepo
	logger *zap.Logger
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=18"`
	Password string `json:"password" binding:"required" validate:"required,min=6,max=18"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
}

type RegisterResponse struct {
	Err    error `json:"err,omitempty"`
	Status int   `json:"status,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=18"`
	Password string `json:"password" binding:"required" validate:"required,min=6,max=18"`
}

type LoginResponse struct {
	Token        string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Err          error  `json:"err,omitempty"`
	Status       int    `json:"status,omitempty"`
}

type GetUserRequest struct {
	ID uint64
}

type GetUserResponse struct {
	User   *models.User `json:"user,omitempty"`
	Err    string       `json:"err,omitempty"`
	Status int          `json:"status,omitempty"`
}

type TaskService interface {
	CreateTask(req CreateTaskRequest) (*CreateTaskResponse, error)
	GetTaskByID(req GetTaskRequest) (*GetTaskResponse, error)
	UpdateTask(req UpdateTaskRequest) (*UpdateTaskResponse, error)
	DeleteTask(req DeleteTaskRequest) (*DeleteTaskResponse, error)
	GetAllTasks(req GetAllTasksRequest) (*GetAllTasksResponse, error)
}

// service implements TaskService
type taskService struct {
	repo   repos.TaskRepo
	logger *zap.Logger
}

type GetTaskRequest struct {
	ID uint64
}

type GetTaskResponse struct {
	Task   *models.Task `json:"task,omitempty"`
	Err    error        `json:"err,omitempty"`
	Status int          `json:"status,omitempty"`
}

type CreateTaskRequest struct {
	Assigned string `json:"assigned" binding:"required" validate:"required,min=3,max=18"`
	Task     string `json:"task" binding:"required" validate:"required,min=3,max=18"`
}

type CreateTaskResponse struct {
	Err    error `json:"err,omitempty"`
	Status int   `json:"status,omitempty"`
}

type UpdateTaskRequest struct {
	ID       uint64 `json:"id"`
	Assigned string `json:"assigned" binding:"required" validate:"required,min=3,max=18"`
	Task     string `json:"task" binding:"required" validate:"required,min=3,max=18"`
}

type UpdateTaskResponse struct {
	Err    error `json:"err,omitempty"`
	Status int   `json:"status,omitempty"`
}

type DeleteTaskRequest struct {
	ID uint64
}

type DeleteTaskResponse struct {
	Err    error `json:"err,omitempty"`
	Status int   `json:"status,omitempty"`
}

type GetAllTasksRequest struct{}

type GetAllTasksResponse struct {
	Task   []*models.Task `json:"tasks,omitempty"`
	Err    error          `json:"err,omitempty"`
	Sum    int            `json:"sum,omitempty"`
	Status int            `json:"status,omitempty"`
}
