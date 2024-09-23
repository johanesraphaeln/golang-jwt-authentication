package handlers

import (
	"golang-jwt-auth/models"
	"golang-jwt-auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Fetch the user from the database
	var storedUser models.User
	if err := models.DB.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a token (assuming the user ID is needed for token generation)
	token, err := utils.GenerateToken(storedUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successfull",
		"data":    token,
	})
}

func Register(c *gin.Context) {
	var user models.User

	// Bind the JSON input to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check for duplicate username
	var storedUser models.User
	if err := models.DB.Where("username = ?", user.Username).First(&storedUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create the new user
	newUser := models.User{
		Username: user.Username,
		Password: string(hashedPassword), // Store the hashed password
	}

	if err := models.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create User"})
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "User Created Successfully",
		"data":    newUser,
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := models.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get All Users Successfully",
		"data":    users,
	})
}
