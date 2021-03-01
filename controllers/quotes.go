package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/carmandomx/quotes/formatter"
	"github.com/carmandomx/quotes/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection = models.ConnectDB()

// CreateQuote saves a Quote to DB
func CreateQuote(c *gin.Context) {

	var quote models.Quote
	id := primitive.NewObjectID()
	if err := c.ShouldBindJSON(&quote); err != nil {
		var verr validator.ValidationErrors

		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": formatter.NewJSONFormatter().Simple(verr)})
			return
		}
		c.Error(err)
		return
	}

	quote.ID = id

	_, err := collection.InsertOne(context.TODO(), quote)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, quote)
}

// GetAllQuotes returns all quotes in collection
func GetAllQuotes(c *gin.Context) {
	options := options.Find()

	options.SetLimit(50)
	cursor, err := collection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		c.Error(err)
		return
	}
	defer cursor.Close(context.TODO())

	var count int
	result := []models.Quote{}
	for cursor.Next(context.TODO()) {
		var helper models.Quote
		err := cursor.Decode(&helper)
		result = append(result, helper)
		if err != nil {
			log.Fatal(err)
		}

		count++

	}
	c.JSON(200, gin.H{
		"count":   count,
		"results": result,
	})

}

// DeleteQuote deletes one from the collection
func DeleteQuote(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	res, err := collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: objectID}})

	if err != nil {
		log.Fatal(err)
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Could not deleted document, try again",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateQuote updates a quote in db
func UpdateQuote(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid ID or could not find document",
		})
		return
	}
	var updatedQuote models.Quote

	if err := c.ShouldBindJSON(&updatedQuote); err != nil {
		var verr validator.ValidationErrors

		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": formatter.NewJSONFormatter().Simple(verr)})
			return
		}
		c.Error(err)
		return
	}
	filter := bson.M{"_id": bson.M{"$eq": objectID}}
	update := bson.M{"$set": bson.M{"quote": updatedQuote.Quote, "class": updatedQuote.Class}}
	after := options.After

	opt := options.FindOneAndUpdate()
	opt.SetReturnDocument(after)
	res := collection.FindOneAndUpdate(context.TODO(), filter, update, opt)

	err = res.Decode(&updatedQuote)

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Error saving to DB, please try again",
		})
		return
	}
	c.JSON(http.StatusOK, updatedQuote)

}
