package ctrls

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/repos"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/services"
	"go.uber.org/zap"
)

func (h *Handler) CreateTask(c *gin.Context) {
	req := &services.CreateTaskRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.env.L().Error("CreateTask request err", zap.Error(err), zap.Reflect("req", c.Request))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	taskRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.L().Error("CreateTask task Repo err: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	taskService := services.NewTaskService(taskRepo, h.env.Log)
	resp, err := taskService.CreateTask(*req)
	if err != nil {
		h.env.L().Error("create task err: ", zap.Error(err))
		c.JSON(resp.Status, gin.H{"error": resp.Err.Error()})
	} else {
		c.JSON(resp.Status, gin.H{"data": resp})
	}
}

func (h *Handler) GetTask(c *gin.Context) {
	req := &services.GetTaskRequest{}
	taskID := c.Param("id")
	req.ID, _ = strconv.ParseUint(taskID, 10, 64)

	taskRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.L().Error("GetTask task Repo err: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	taskService := services.NewTaskService(taskRepo, h.env.Log)
	resp, _ := taskService.GetTaskByID(*req)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) UpdateTask(c *gin.Context) {
	req := &services.UpdateTaskRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.env.L().Error("UpdateTask request err", zap.Error(err), zap.Reflect("req", c.Request))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskID := c.Param("id")
	req.ID, _ = strconv.ParseUint(taskID, 10, 64)

	taskRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.L().Error("UpdateTask task Repo err: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	taskService := services.NewTaskService(taskRepo, h.env.Log)
	resp, _ := taskService.UpdateTask(*req)

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	req := &services.DeleteTaskRequest{}
	taskID := c.Param("id")
	req.ID, _ = strconv.ParseUint(taskID, 10, 64)

	taskRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.L().Error("DeleteTask task Repo err: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	taskService := services.NewTaskService(taskRepo, h.env.Log)
	resp, _ := taskService.DeleteTask(*req)

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	req := &services.GetAllTasksRequest{}
	taskRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.L().Error("GetAllTasks task Repo err: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	taskService := services.NewTaskService(taskRepo, h.env.Log)
	resp, _ := taskService.GetAllTasks(*req)

	c.JSON(http.StatusOK, gin.H{"data": resp})
}
