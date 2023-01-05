package ctrls

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/repos"
)

func GetUser(c *gin.Context) {
	repo := repos.NewUserRepo()                       // Create user repo
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64) // Parse ID from URL
	req := repos.GetUserRequest{ID: id}

	if resp, err := repo.GetUser(req); err != nil { // Try to get user from database
		c.AbortWithStatus(http.StatusNotFound) // Abort if not found
		log.Println(err)
	} else {
		c.JSON(http.StatusOK, resp) // Send back data
	}
}
