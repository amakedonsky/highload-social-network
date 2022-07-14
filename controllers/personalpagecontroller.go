package controllers

import (
	"amakedonsky/highload-social-network/database"
	"amakedonsky/highload-social-network/helpers"
	"amakedonsky/highload-social-network/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePersonalPage(c *gin.Context) {
	var newPage *models.PersonalPage

	if err := c.BindJSON(&newPage); err != nil {
		c.Error(err)
		return
	}

	id, err := database.CreatePersonalPage(c.Request.Context(), *newPage)
	if err != nil {
		c.Error(err)
		return
	}
	newPage.Id = strconv.FormatInt(id, 10)

	c.IndentedJSON(http.StatusCreated, newPage)
}

func UpdatePersonalPage(c *gin.Context) {
	user := c.MustGet("User").(models.PersonalPage)

	if user.Email == "" {
		helpers.ResponseWithError(c, http.StatusForbidden, "Please login first")
		return
	}

	var newPage *models.PersonalPage
	if err := c.BindJSON(&newPage); err != nil {
		c.Error(err)
		return
	}
	newPage.Id = user.Id

	_, err := database.UpdatePersonalPage(c.Request.Context(), *newPage)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, newPage)
}

func GetPersonalPage(c *gin.Context) {
	user := c.MustGet("User").(models.PersonalPage)

	if user.Email == "" {
		helpers.ResponseWithError(c, http.StatusForbidden, "Please login first")
		return
	}

	result, err := database.FetchPersonalPage(c.Request.Context(), user.Id)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
