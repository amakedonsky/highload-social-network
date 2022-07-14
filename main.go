package main

import (
	"amakedonsky/highload-social-network/database"
	"amakedonsky/highload-social-network/middleware"
	"amakedonsky/highload-social-network/routes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()

	log.Println("Starting the HTTP server on port 8090")
	router := gin.Default()
	initHandlers(router)
	log.Fatal(http.ListenAndServe(":8090", router))
}

func initHandlers(router *gin.Engine) {
	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/api/v1")
	private.Use(middleware.Authenticate())
	routes.PrivateRoutes(private)
}
