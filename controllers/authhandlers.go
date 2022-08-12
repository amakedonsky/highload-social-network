package controllers

import (
	"amakedonsky/highload-social-network/database"
	"amakedonsky/highload-social-network/helpers"
	"amakedonsky/highload-social-network/models"
	"amakedonsky/highload-social-network/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.SignInReq

		// Bind the request body data to var data and check if all details are provided
		if c.BindJSON(&data) != nil {
			c.JSON(406, gin.H{"message": "Provide required details"})
			c.Abort()
			return
		}

		result, err := database.GetPasswordByEmail(c.Request.Context(), data.Email)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Problem logging into your account"})
			c.Abort()
			return
		}

		if result.Email == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": "User account was not found"})
			c.Abort()
			return
		}

		// Get the hashed password from the saved document
		hashedPassword := []byte(result.Password)
		// Get the password provided in the request.body
		password := []byte(data.Password)

		err = helpers.PasswordCompare(password, hashedPassword)

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid user credentials"})
			c.Abort()
			return
		}

		jwtToken, _, err2 := services.GenerateToken(data.Email)

		// If we fail to generate token for access
		if err2 != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "There was a problem logging you in, try again later"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Log in success", "token": jwtToken})
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newPage *models.PersonalPage

		if err := c.BindJSON(&newPage); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": err})
			c.Abort()
			return
		}

		if helpers.EmptyUserPass(newPage.Email, newPage.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Empty password or email"})
			c.Abort()
			return
		}

		_, err := database.CreatePersonalPage(c.Request.Context(), *newPage)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			c.Abort()
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "New user page registered"})
	}
}
