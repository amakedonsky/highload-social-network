package routes

import (
	"amakedonsky/highload-social-network/controllers"
	"github.com/gin-gonic/gin"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.POST("/signin", controllers.Signin())
	g.POST("/signup", controllers.Signup())
}

func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/page", controllers.GetPersonalPage)
	g.POST("/page", controllers.CreatePersonalPage)
	g.PUT("/page", controllers.UpdatePersonalPage)

	g.POST("/friends/:id", controllers.AddToFriends)
	g.DELETE("/friends/:id", controllers.DelFromFriends)
	g.GET("/friends", controllers.GetAllFriends)
}
