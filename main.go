package main

import (
	"log"
	"day4/controllers"
	"day4/database"
	"day4/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No env file found")
	}

	database.ConnectDB()

	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/login", controllers.Login)
	protected := r.Group("/api")
	protected.Use(middleware.Middleware())
	protected.GET("/dashboard", controllers.Dashboard)
	protected.GET("/users", controllers.GetUser)
	protected.POST("/logout", controllers.Logout)
	r.Run(":8080")

}