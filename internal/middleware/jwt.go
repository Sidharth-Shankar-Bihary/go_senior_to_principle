package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramseyjiang/go_senior_to_principle/internal/env"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/auth"
)

// Auth returns a Gin HandlerFunc function.
// This function expects a context for which it tries to validate the JWT in the header.
// If it is invalid, an error response is returned. If not, the Next() function on the context is called.
func Auth(e *env.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenDetail, err := auth.ExtractTokenData(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
		}

		userID, err := auth.FetchUIDFromRedis(accessTokenDetail, e.Redis)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
		}

		if userID != 0 {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthenticated",
			})
			c.Abort()
			return
		}
	}
}
