package controllers

import (
	"day4/database"
	"day4/models"
	"day4/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "hashin password fails",
		})
		return
	}
	user.Password = hash
	if user.Role == "" {
		user.Role = "user"
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(401, gin.H{
			"error": "invalid input",
		})
		return
	}
	if err := database.DB.Where("email=?", input.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{
			"error": "wrong password",
		})
		return
	}
	if err := utils.ComparePassword(user.Password, input.Password); err != nil {
		c.JSON(401, gin.H{"error": "wrong password"})
		return
	}
	access, err := utils.GenarateAccessToken(user.Id, user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to access token"})
		return
	}
	refresh, err := utils.GenarateRefreshToken(user.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to genarate refresh token"})
		return
	}
	user.RefreshToken = refresh
	database.DB.Save(&user)
	c.JSON(200, gin.H{
		"access":  access,
		"refresh": refresh,
	})
}

func Dashboard(c *gin.Context) {
	role := c.GetString("role")
	if role == "admin" {
		c.JSON(200, gin.H{
			"message": "welcome to admin page",
		})
		return
	}
	c.JSON(200, gin.H{"message": "welcome to user"})
}

func GetUser(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(403, gin.H{
			"error": "access denied",
		})
		return
	}
	var users []models.User
	if err := database.DB.Select("id,name,email,role").Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch users"})
		return
	}
	c.JSON(200, users)
}

func Logout(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := database.DB.Model(&models.User{}).
		Where("id=?", userID).Update("refresh_token", "").Error; err != nil {
		c.JSON(500, gin.H{"error": "logout failed"})
		return
	}
	c.JSON(200, gin.H{"message": "logged out sucessfully"})
}