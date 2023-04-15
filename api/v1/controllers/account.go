package controllers

import (
	"context"
	"net/http"

	"saasmanagement/api/v1/validators"
	"saasmanagement/config"
	"saasmanagement/models"
	"saasmanagement/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAccount(c *gin.Context) {
	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the user input
	if err := validators.ValidateAccount(&account); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Insert the user into the database
	result, err := config.GetDBCollection("accounts").InsertOne(context.Background(), account)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Return the created account id
	utils.Success(c, result.InsertedID)
}

func UpdateAccount(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ObjectID")
		return
	}

	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validators.ValidateAccount(&account); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": account}

	result, err := config.GetDBCollection("accounts").UpdateOne(context.Background(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		utils.Error(c, http.StatusInternalServerError, "Failed to update account")
		return
	}

	utils.Success(c, "Account updated")
}

func DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ObjectID")
		return
	}

	filter := bson.M{"_id": objectID}
	result, err := config.GetDBCollection("accounts").DeleteOne(context.Background(), filter)
	if err != nil || result.DeletedCount == 0 {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	utils.Success(c, "Account deleted")
}

func GetAllAccounts(c *gin.Context) {
	var accounts []models.Account
	cursor, err := config.GetDBCollection("accounts").Find(context.Background(), bson.M{})
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get accounts")
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var account models.Account
		cursor.Decode(&account)
		accounts = append(accounts, account)
	}

	utils.Success(c, accounts)
}

func GetAccountByID(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ObjectID")
		return
	}

	filter := bson.M{"_id": objectID}
	var account models.Account
	err = config.GetDBCollection("accounts").FindOne(context.Background(), filter).Decode(&account)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get account")
		return
	}

	utils.Success(c, account)
}
