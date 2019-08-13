package db

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var mongoHost = os.Getenv("APP_MONGO_HOST")
var mongoDatabaseName = os.Getenv("APP_MONGO_DATABASE")

func ConnectMiddleware() func(*gin.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoHost))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("ping %s:%s ping error: %v\n", err)
		panic(err)
	}
	database := client.Database(mongoDatabaseName)

	if database == nil {
		panic("Database is not connected!")
	} else {
		fmt.Println("Database is connected.")
	}

	return func(c *gin.Context) {
		c.Set("database", database)
		c.Next()
	}
}
