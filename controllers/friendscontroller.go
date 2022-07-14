package controllers

import (
	"amakedonsky/highload-social-network/database"
	"amakedonsky/highload-social-network/helpers"
	"amakedonsky/highload-social-network/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddToFriends(c *gin.Context) {
	friendId := c.Param("id")

	user := c.MustGet("User").(models.PersonalPage)

	if user.Email == "" {
		helpers.ResponseWithError(c, http.StatusForbidden, "Please login first")
		return
	}

	err := database.AddToFriends(c.Request.Context(), user.Id, friendId)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, "Added")
}

func DelFromFriends(c *gin.Context) {
	friendId := c.Param("id")

	user := c.MustGet("User").(models.PersonalPage)

	if user.Email == "" {
		helpers.ResponseWithError(c, http.StatusForbidden, "Please login first")
		return
	}

	err := database.DelFromFriends(c.Request.Context(), user.Id, friendId)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, "Deleted")
}

func GetAllFriends(c *gin.Context) {
	user := c.MustGet("User").(models.PersonalPage)

	if user.Email == "" {
		helpers.ResponseWithError(c, http.StatusForbidden, "Please login first")
		return
	}

	pages, err := database.GetAllFriends(c.Request.Context(), user.Id)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, pages)
}
