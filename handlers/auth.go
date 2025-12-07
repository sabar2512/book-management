package handlers

import (
	"book-management/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user authentication and token generation
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Simple authentication (in production, verify against database with hashed passwords)
	// For demo purposes, accepting any username/password combination
	// You should implement proper authentication here
	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"username": input.Username,
	})
}
