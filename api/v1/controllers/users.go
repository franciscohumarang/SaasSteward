package controllers

import (
	"context"
	"log"
	"net/http"

	"saasmanagement/api/v1/validators"
	"saasmanagement/config"
	"saasmanagement/models"
	"saasmanagement/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the user input
	if err := validators.ValidateUser(&user); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if the user already exists
	existingUser := &models.User{}
	err := config.GetDBCollection("users").FindOne(context.Background(), bson.M{"email": user.Email}).Decode(existingUser)
	if err == nil {
		utils.Error(c, http.StatusConflict, "User with the same email already exists")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hashedPassword)

	// Insert the user into the database
	result, err := config.GetDBCollection("users").InsertOne(context.Background(), user)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Return the created user ID
	utils.Success(c, result.InsertedID)
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	cursor, err := config.GetDBCollection("users").Find(context.Background(), bson.M{})
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get users")
		return
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	utils.Success(c, users)
}

func GetUserById(c *gin.Context) {
	userId := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	err = config.GetDBCollection("users").FindOne(context.Background(), bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, user)
}

func UpdateUser(c *gin.Context) {
	userId := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the user input
	if err := validators.ValidateUser(&user); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if the user already exists
	existingUser := &models.User{}
	err = config.GetDBCollection("users").FindOne(context.Background(), bson.M{"email": user.Email}).Decode(existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found, allow update
		} else {
			utils.Error(c, http.StatusInternalServerError, "Failed to update user")
			return
		}
	} else if existingUser.ID.Hex() != userId {
		utils.Error(c, http.StatusConflict, "User with the same email already exists")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hashedPassword)

	// Update the user in the database
	update := bson.M{
		"$set": user,
	}
	result, err := config.GetDBCollection("users").UpdateOne(context.Background(), bson.M{"_id": objId}, update)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	if result.ModifiedCount == 0 {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, result.ModifiedCount)
}

func DeleteUser(c *gin.Context) {
	userId := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	result, err := config.GetDBCollection("users").DeleteOne(context.Background(), bson.M{"_id": objId})
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	if result.DeletedCount == 0 {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, result.DeletedCount)
}

func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the input data
	if err := validators.ValidateUser(&models.User{Email: loginData.Email, Password: loginData.Password}); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if the user exists
	var user models.User
	err := config.GetDBCollection("users").FindOne(context.Background(), bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.Error(c, http.StatusUnauthorized, "Incorrect email or password")
		} else {
			utils.Error(c, http.StatusInternalServerError, "Failed to login")
		}
		return
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Generate a JWT token
	token, err := config.GenerateToken(user.ID.Hex())
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{"token": token})
}
