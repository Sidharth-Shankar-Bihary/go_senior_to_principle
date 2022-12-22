package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/ctrls"
	"net/http"
)

type Router struct {
	router *gin.Engine
}

func (r *Router) Handler() http.Handler {
	return r.router.Handler()
}

func path(endpoint string) string {
	return fmt.Sprintf("/api/%s", endpoint)
}

func NewRouter() *Router {
	router := gin.Default()

	// use cors middleware. cors.Default() means allow all origins. If replace cors.Default() with cors.New(), you define allow by yourself.
	router.Use(cors.Default())

	router.Handle("GET", path("health"), ctrls.CheckHealth)
	router.Handle("GET", path("users/:id"), ctrls.GetUser)

	return &Router{
		router,
	}
}

// func (r *Router) GetProjects(c *gin.Context) {
// 	projects, err := services.GetProjectModel().GetAllProjects()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"data": projects,
// 	})
// }
//
// func (r *Router) EditProject(c *gin.Context) {
// 	id := c.Param("id")
// 	var project models.Project
// 	if err := c.ShouldBindJSON(&project); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	idInt, err := strconv.ParseUint(id, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id, not a number"})
// 		return
// 	}
// 	if err := services.GetProjectModel().EditProject(uint(idInt), &project); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"id": project.ID,
// 	})
// }
//
// func (r *Router) AddProject(c *gin.Context) {
// 	var project models.Project
// 	if err := c.ShouldBindJSON(&project); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := services.GetProjectModel().CreateProject(&project); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add project"})
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"id": project.ID,
// 	})
// }
//
// func (r *Router) DeleteProject(c *gin.Context) {
// 	id := c.Param("id")
// 	idInt, err := strconv.ParseUint(id, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id, not a number"})
// 		return
// 	}
//
// 	if err := services.GetProjectModel().DeleteProject(uint(idInt)); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"id": id,
// 	})
// }
