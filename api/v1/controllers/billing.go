package controllers

import (
	"context"
	"net/http"
	"saasmanagement/api/v1/validators"
	"saasmanagement/config"
	"saasmanagement/models"
	"saasmanagement/utils"

	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get a single billing by ID
func GetBillingByID(c *gin.Context) {

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	filter := bson.M{"_id": id}
	var billing models.Billing
	err = config.GetDBCollection("billings").FindOne(context.Background(), filter).Decode(&billing)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Billing not found"})
		return
	}

	utils.Success(c, billing)
}

func CreateBilling(c *gin.Context) {
	var billing models.Billing
	if err := c.ShouldBindJSON(&billing); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the user input
	if err := validators.ValidateBilling(&billing); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Insert the user into the database
	result, err := config.GetDBCollection("billings").InsertOne(context.Background(), billing)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create billing")
		return
	}

	// Return the created account id
	utils.Success(c, result.InsertedID)
}

func UpdateBilling(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var billing models.Billing
	if err := c.ShouldBindJSON(&billing); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validators.ValidateBilling(&billing); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": billing}

	result, err := config.GetDBCollection("billings").UpdateOne(context.Background(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		utils.Error(c, http.StatusInternalServerError, "Failed to update billing")
		return
	}

	utils.Success(c, "Account updated")
}

func GetAllBillingsByUserID(c *gin.Context) {
	objectID := c.Query("user_id")

	// Get pagination parameters from GET query parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Convert parameters to integers
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		limitNum = 10
	}

	// Calculate skip value for pagination
	skip := (pageNum - 1) * limitNum

	// Set up filter and options for database query
	filter := bson.M{"user_id": objectID}
	options := options.Find()
	options.SetSkip(int64(skip))
	options.SetLimit(int64(limitNum))
	options.SetSort(bson.M{"date_created": -1})

	// Execute database query and handle errors
	cursor, err := config.GetDBCollection("billings").Find(context.Background(), filter, options)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get billings from database")
		return
	}
	defer cursor.Close(context.Background())

	var billings []models.Billing
	// Iterate over results and append to billings slice
	for cursor.Next(context.Background()) {
		var billing models.Billing
		cursor.Decode(&billing)
		billings = append(billings, billing)
	}

	// Return paginated billings in response
	utils.Success(c, billings)
}

func DeleteBilling(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ObjectID")
		return
	}

	filter := bson.M{"_id": objectID}
	result, err := config.GetDBCollection("billings").DeleteOne(context.Background(), filter)
	if err != nil || result.DeletedCount == 0 {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete billing")
		return
	}

	utils.Success(c, "Billing deleted")
}
