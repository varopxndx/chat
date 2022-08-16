package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthToken middleware to use a token as user session
func AuthToken(validateToken func(string) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("authorization")
		var token string
		if value == "" {
			token = c.Query("bearer")
		} else {
			words := strings.Fields(value)

			// validate if value is empty
			if len(words) != 2 {
				token = ""
			} else {
				token = words[1]
			}
		}
		err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
