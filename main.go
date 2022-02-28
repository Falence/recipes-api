// swagger:meta
//
// Recipe API
//
// This is a sample recipes API. Youcan find out more about the API at
// https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Falence Lemungoh
// <falencelemungoh@gmail.com>https://github.com/falence
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	// "strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipe struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Tags         []string           `bson:"tags" json:"tags"`
	Ingredients  []string           `bson:"ingredients" json:"ingredients"`
	Instructions []string           `bson:"instructions" json:"instructions"`
	PublishedAt  time.Time          `bson:"publishedAt" json:"publishedAt"`
}

var (
	ctx        context.Context
	err        error
	client     *mongo.Client
	collection *mongo.Collection
)

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("err")
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// recipe.ID = xid.New().String()
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting a new recipe",
		})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer cur.Close(ctx)
	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "name", Value: recipe.Name},
				primitive.E{Key: "instructions", Value: recipe.Instructions},
				primitive.E{Key: "ingredients", Value: recipe.Ingredients},
				primitive.E{Key: "tags", Value: recipe.Tags},
			}},
		},
	)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been updated",
	})
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	res, err := collection.DeleteOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: objectId},
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

func SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	cur, err := collection.Find(ctx, bson.M{"tags": tag})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer cur.Close(ctx)
	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.Run()
}
