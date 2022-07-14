package helpers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func GeneratePasswordHash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}

func PasswordCompare(password []byte, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)

	return err
}

func ResponseWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"message": message})
}
