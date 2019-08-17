package controllers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func badRequest(c *gin.Context, errorCode string) {
	c.JSON(400, gin.H{
		"error": errorCode,
	})
}

func internalError(c *gin.Context, errorCode string) {
	c.JSON(500, gin.H{
		"error": errorCode,
	})
}

func TypeformController(r *gin.Engine) {
	typeformGroup := r.Group("/typeform")

	typeformGroup.POST("/webhook", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			badRequest(c, "INVALID_REQUEST")
		}
		var payload interface{}
		err = json.Unmarshal(body, &payload)
		if err != nil {
			badRequest(c, "INVALID_JSON")
		}

		db := c.MustGet("database").(*mongo.Database)
		responses := db.Collection("responses")

		ctx := context.Background()

		_, err = responses.InsertOne(ctx, payload)
		if err != nil {
			internalError(c, "DATABASE_INSERT_ERROR")
		}

		c.JSON(200, gin.H{})
	})

	antiUsers := map[string]string{
		"10213783272831334": "Danielle Felker",
	}

	typeformGroup.GET("/survey", func(c *gin.Context) {
		fbId := c.Query("id")
		location := "https://vinlock1.typeform.com/to/CWy6cX?name=" + antiUsers[fbId] + "&id=" + fbId
		c.Redirect(302, location)
	})
}
