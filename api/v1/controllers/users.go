package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"saasmanagement/api/v1/validators"
	"saasmanagement/config"
	"saasmanagement/models"
	"saasmanagement/utils"
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
	var loginData models.Login

	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the input data
	if err := validators.ValidateLogin(&loginData); err != nil {
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
	/*token, err := config.GenerateToken(user.ID.Hex())
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	*/
	// if authentication succeeds, generate an access token and a refresh token
	accessClaims := models.Claims{
		UserId:       user.ID.Hex(),
		AccessToken:  true,
		RefreshToken: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    "saasteward",
		},
	}
	var jwtSecret = []byte(config.GetEnv("JWT_SECRET"))
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString(jwtSecret)

	if err != nil {

		utils.Error(c, http.StatusInternalServerError, "Failed to generate signed access token")
		return
	}

	refreshClaims := models.Claims{
		UserId:       user.ID.Hex(),
		AccessToken:  false,
		RefreshToken: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
			Issuer:    "saasteward",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate signed refresh token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  signedAccessToken,
		"refreshToken": signedRefreshToken,
	})

}

// RefreshHandler generates a new access token using a valid refresh token
func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	var jwtSecret = []byte(config.GetEnv("JWT_SECRET"))
	if err != nil {

		utils.Error(c, http.StatusBadRequest, "Refresh token missing")
		return
	}

	token, err := jwt.ParseWithClaims(refreshToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {

		utils.Error(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !claims.RefreshToken || claims.ExpiresAt < time.Now().Unix() {
		utils.Error(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	accessClaims := models.Claims{
		UserId:       claims.UserId,
		AccessToken:  true,
		RefreshToken: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    "your-app-name",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": signedAccessToken})
}
