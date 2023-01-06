package ctrls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramseyjiang/go_senior_to_principle/internal/env"
)

type Handler struct {
	env *env.Environment
}

func NewHandler(hEnv *env.Environment) *Handler {
	return &Handler{
		env: hEnv,
	}
}

func (h *Handler) CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "Health is ok."})
}
