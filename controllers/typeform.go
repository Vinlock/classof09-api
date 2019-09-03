package controllers

import (
	"context"
	"ecr-reunion/typeform"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
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
		"Y3J5c3RhbC5hdHRhbGxhQGdtYWlsLmNvbQo=": "Crystal Attalla",
		"Y2luZHl5X3JlbGxhMDlAeWFob28uY29tCg==": "Cindy Caballero",
	}

	typeformGroup.GET("/survey", func(c *gin.Context) {
		valid := true
		fbId, ok := c.GetQuery("id")
		name, nameOk := antiUsers[fbId]
		if ok && nameOk {
			typeformApi := typeform.NewTypeformApi(os.Getenv("APP_TYPEFORM_TOKEN"))
			completedOnly := true
			params := typeform.GetResponsesParams{
				FormId:    os.Getenv("APP_SURVEY1_ID"),
				Completed: &completedOnly,
				Query:     fbId,
			}
			response, err := typeformApi.GetResponses(params)
			if err != nil || response.TotalItems > 0 {
				valid = false
			} else {
				surveyId := os.Getenv("APP_SURVEY1_ID")
				location := "https://vinlock1.typeform.com/to/" + surveyId + "?name=" + name + "&id=" + fbId
				c.Redirect(302, location)
			}
		} else {
			valid = false
		}

		if !valid {
			c.Redirect(302, "https://classof09.org")
		}
	})
}
