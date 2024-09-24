package handlers

import (
	"context"
	"fmt"
	"golang-jwt-auth/models"
	"golang-jwt-auth/utils"
	"log"
	"net/http"

	"firebase.google.com/go/v4/db"
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
		"message": "User Created Successfully in PostgreSQL",
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

type Project struct {
	Name string `json:"name"`
}

func SaveProject(c *gin.Context, client *db.Client) {
	var project Project

	// Bind the JSON input to the project struct
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Retrieve the user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// create ref at path user_projects/:userId
	ref := client.NewRef("user_projects/" + fmt.Sprint(userID))

	if err := ref.Set(context.TODO(), map[string]interface{}{
		"name": project.Name,
	}); err != nil {
		log.Fatal(err)
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "Project Saved Successfully in Firebase",
		"data":    project,
	})
}

func GetProject(c *gin.Context, client *db.Client) {
	// Retrieve the user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// get database reference to user project
	ref := client.NewRef("user_projects/" + fmt.Sprint(userID))

	// read from user_projects using ref
	var project Project
	if err := ref.Get(context.TODO(), &project); err != nil {
		log.Fatalln("error in reading from firebase DB: ", err)
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Project Retrieved Successfully in Firebase",
		"data":    project,
	})
}

func DeleteProject(c *gin.Context, client *db.Client) {
	// Retrieve the user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// get database reference to user project
	ref := client.NewRef("user_projects/" + fmt.Sprint(userID))

	if err := ref.Delete(context.TODO()); err != nil {
		log.Fatalln("error in deleting ref: ", err)
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Project Deleted Successfully in Firebase",
		"data":    nil,
	})
}
