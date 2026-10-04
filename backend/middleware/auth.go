package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		expectedKey := os.Getenv("APPROVAL_API_KEY")

		if expectedKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "approval API key is not configured",
			})
			c.Abort()
			return
		}

		providedKey := c.GetHeader("X-API-Key")

		if providedKey == "" || providedKey != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Set("reviewer", "api-reviewer")
		c.Next()
	}
}