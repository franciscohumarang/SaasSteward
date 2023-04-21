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
func GetCompanyByID(c *gin.Context) {

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	filter := bson.M{"_id": id}
	var company models.Company
	err = config.GetDBCollection("companies").FindOne(context.Background(), filter).Decode(&company)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	utils.Success(c, company)
}

func CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the user input
	if err := validators.ValidateCompany(&company); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Insert the user into the database
	result, err := config.GetDBCollection("companies").InsertOne(context.Background(), company)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create billing")
		return
	}

	// Return the created account id
	utils.Success(c, result.InsertedID)
}

func UpdateCompany(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validators.ValidateCompany(&company); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": company}

	result, err := config.GetDBCollection("billings").UpdateOne(context.Background(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		utils.Error(c, http.StatusInternalServerError, "Failed to update company")
		return
	}

	utils.Success(c, "Company updated")
}

func GetAllCompanysByUserID(c *gin.Context) {
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

	var billings []models.Company
	// Iterate over results and append to billings slice
	for cursor.Next(context.Background()) {
		var billing models.Company
		cursor.Decode(&billing)
		billings = append(billings, billing)
	}

	// Return paginated billings in response
	utils.Success(c, billings)
}

func DeleteCompany(c *gin.Context) {
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

	utils.Success(c, "Company deleted")
}
