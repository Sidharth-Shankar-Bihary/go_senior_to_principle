package services

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/repos"
	"go.uber.org/zap"
)

func NewTaskService(repo repos.TaskRepo, logger *zap.Logger) TaskService {
	return &taskService{
		repo:   repo,
		logger: logger,
	}
}

func (t *taskService) CreateTask(req CreateTaskRequest) (resp *CreateTaskResponse, err error) {
	// give the original value to resp, otherwise it will make invalid memory address or nil pointer dereference
	resp = &CreateTaskResponse{}
	task := models.Task{}

	// valid request
	if err = validator.New().Struct(req); err != nil {
		resp.Err = err
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	task.Assigned = req.Assigned
	task.Task = req.Task
	if err = t.repo.CreateTask(&task); err != nil {
		resp.Err = err
		resp.Status = http.StatusBadRequest
		return
	}

	resp.Status = http.StatusOK
	resp.Err = nil
	return
}

func (t *taskService) GetTaskByID(req GetTaskRequest) (resp *GetTaskResponse, err error) {
	resp = &GetTaskResponse{}
	if req.ID <= 0 {
		resp.Status = http.StatusBadRequest
		resp.Err = errors.New("request ID is wrong")
		return
	}

	task, err := t.repo.GetTaskByID(uint(req.ID))
	if err != nil {
		resp.Err = err
		resp.Status = http.StatusBadRequest
		return
	}

	if task.Task == "" && task.Assigned == "" {
		resp.Task = nil
		resp.Status = http.StatusNotFound
		resp.Err = nil
		return
	}
	resp.Task = task
	resp.Status = http.StatusOK
	resp.Err = nil

	return
}

func (t *taskService) UpdateTask(req UpdateTaskRequest) (resp *UpdateTaskResponse, err error) {
	resp = &UpdateTaskResponse{}
	task := models.Task{}

	if err = validator.New().Struct(req); err != nil {
		resp.Err = err
		resp.Status = http.StatusUnprocessableEntity
		return
	}

	task.Assigned = req.Assigned
	task.Task = req.Task
	task.ID = uint(req.ID)
	if err = t.repo.UpdateTask(&task); err != nil {
		resp.Err = err
		resp.Status = http.StatusBadRequest
		return
	}

	resp.Status = http.StatusOK
	resp.Err = nil
	return
}
func (t *taskService) DeleteTask(req DeleteTaskRequest) (resp *DeleteTaskResponse, err error) {
	resp = &DeleteTaskResponse{}
	if req.ID <= 0 {
		resp.Status = http.StatusBadRequest
		resp.Err = errors.New("request ID is wrong")
		return
	}

	if err = t.repo.DeleteTask(uint(req.ID)); err != nil {
		resp.Err = err
		resp.Status = http.StatusBadRequest
		return
	}

	resp.Status = http.StatusOK
	resp.Err = nil
	return
}

func (t *taskService) GetAllTasks(_ GetAllTasksRequest) (resp *GetAllTasksResponse, err error) {
	resp = &GetAllTasksResponse{}
	tasks, err := t.repo.GetAllTasks()
	if err != nil {
		resp.Task = nil
		resp.Status = http.StatusBadRequest
		resp.Err = errors.New("get all tasks failed")
	}

	resp.Sum = len(tasks)
	resp.Task = tasks
	resp.Status = http.StatusOK
	resp.Err = nil
	return
}
