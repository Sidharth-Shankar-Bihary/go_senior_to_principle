package ctrls

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/repos"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/services"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/auth"
	"go.uber.org/zap"
)

func (h *Handler) Register(c *gin.Context) {
	req := &services.RegisterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.env.L().Error("Register request err", zap.Error(err), zap.Reflect("req", c.Request))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.L().Error("Register user Repo err: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userService := services.NewUserService(userRepo, h.env.Log)
	resp, err := userService.CreateUser(*req)
	if err != nil {
		h.env.L().Error("Register create user err: ", zap.Error(err))
		c.JSON(resp.Status, gin.H{"error": resp.Err.Error()})
	} else {
		c.JSON(resp.Status, gin.H{"data": resp})
	}
}

func (h *Handler) Login(c *gin.Context) {
	req := &services.LoginRequest{}
	if err := c.ShouldBind(&req); err != nil {
		h.env.L().Error("Login request err: ", zap.Error(err), zap.Reflect("req", c.Request))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userRepo, err := repos.New(h.env.DB, h.env.Log)
	if err != nil {
		h.env.L().Error("Login user Repo err: ", zap.Error(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	userService := services.NewUserService(userRepo, h.env.Log)
	resp, err := userService.GetUserToken(h.env.Redis, *req)

	if err != nil {
		h.env.L().Error("Login user service err: ", zap.Error(err))
		c.JSON(resp.Status, gin.H{"error": resp.Err.Error()})
	} else {
		// set jwt token in a cookie for frontend use it.
		http.SetCookie(c.Writer, &http.Cookie{
			Name:   "user_token",
			Value:  resp.Token,
			MaxAge: 0,
			Path:   "/",
		})

		c.JSON(resp.Status, gin.H{"data": resp})
	}
}

// Logout needs to clear redis, when generate token should be saved it in redis.
func (h *Handler) Logout(c *gin.Context) {
	accessTokenDetail, err := auth.ExtractTokenData(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	err = auth.DelAuth(accessTokenDetail.AccessUUID, h.env.Redis)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// set cookie expires is used for frontend.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "user_token",
		Value:   "",
		MaxAge:  0,
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Logged out successfully",
	})
}

// RefreshToken uses the refresh_token to refresh new access token and new refresh token, after that the old refresh token is invalid.
func (h *Handler) RefreshToken(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	token, err := auth.ParseRefreshToken(mapToken["refresh_token"])
	if err != nil { // if there is an error, the token must have expired
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}

	// check refresh token is valid or not
	if _, ok := token.Claims.(jwt.RegisteredClaims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	tokens, err := auth.RegenerateToken(token, h.env.Redis)
	if err != nil {
		h.env.L().Error("RefreshToken route RegenerateToken err: ", zap.Error(err))
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tokens})
}
