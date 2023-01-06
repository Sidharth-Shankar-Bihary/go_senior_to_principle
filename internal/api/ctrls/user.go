package ctrls

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/repos"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/services"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/auth"
	"go.uber.org/zap"
)

func (h *Handler) GetUser(c *gin.Context) {
	req := &services.GetUserRequest{}
	userID := c.Param("id")
	req.ID, _ = strconv.ParseUint(userID, 10, 64)

	userRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.Log.Error("GetUser init userRepo err: ", zap.Any("error", err))
		os.Exit(0)
	}
	userService := services.NewUserService(userRepo, h.env.Log)
	resp, _ := userService.GetUserByID(*req)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) CurrentUser(c *gin.Context) {
	accessTokenDetail, err := auth.ExtractTokenData(c)
	if err != nil {
		h.env.L().Error("CurrentUser route accessTokenDetail err: ", zap.Error(err), zap.Reflect("req: ", c.Request))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	userID, err := auth.FetchUIDFromRedis(accessTokenDetail, h.env.Redis)
	if err != nil {
		h.env.Log.Error("CurrentUser FetchUIDFromRedis err: ", zap.Any("error", err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	req := &services.GetUserRequest{}
	req.ID = userID

	userRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.Log.Error("CurrentUser route userRepo err: ", zap.Any("error", err))
		os.Exit(0)
	}
	userService := services.NewUserService(userRepo, h.env.Log)
	resp, _ := userService.GetUserByID(*req)
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
