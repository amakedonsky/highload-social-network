package middleware

import (
	"amakedonsky/highload-social-network/database"
	"amakedonsky/highload-social-network/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func responseWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"message": message})
}

// Authenticate fetches user details from token
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		authToken := c.Request.Header["Authorization"]

		if len(authToken) == 0 {
			responseWithError(c, http.StatusForbidden, "Please login to your account")
			return
		}

		words := strings.Fields(authToken[0])
		username, _ := services.DecodeToken(words[1])

		if username == "" {
			responseWithError(c, http.StatusNotFound, "User account not found")
			return
		}

		result, err := database.GetUserByEmail(c.Request.Context(), username)

		if result.Email == "" {
			responseWithError(c, http.StatusNotFound, "User account not found")
			return
		}

		if err != nil {
			responseWithError(c, http.StatusInternalServerError, "Something wrong with your account")
			return
		}

		c.Set("User", result)
		c.Next()
	}
}
