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
	"os"

	"github.com/falence/recipes-api/handlers"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	authHandler    *handlers.AuthHandler
	recipesHandler *handlers.RecipesHandler
)

func init() {
	ctx := context.Background()
	mongoOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("err")
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)
	recipesHandler = handlers.NewRecipeHandler(ctx, collection, redisClient)
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	store, _ := redisStore.NewStore(10, "tcp", os.Getenv("REDIS_URI"), "", []byte(os.Getenv("REDIS_SESSION_SECRET")))

	router := gin.Default()
	router.Use(sessions.Sessions("recipes_api", store))

	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/signup", authHandler.SignupHandler)
	router.POST("/signout", authHandler.SignOutHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.GET("/recipes/:id", recipesHandler.FindRecipeHandler)
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.GET("/recipes/search", recipesHandler.SearchRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}
	router.Run()
}
