package main

import (
	"ecr-reunion/auth"
	"ecr-reunion/db"
	"ecr-reunion/typeform"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func main() {
	ginIsReleaseMode := os.Getenv("APP_GIN_MODE_RELEASE") == "true"
	if ginIsReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Cors
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST"},
		AllowHeaders:    []string{"Origin", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))

	r.Use(db.ConnectMiddleware())

	// Typeform Controller
	typeform.TypeformController(r)

	// JWT Middleware
	auth.JWTMiddleware(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":  "PAGE_NOT_FOUND",
			"error": "Page not found",
		})
	})

	appPort := os.Getenv("APP_PORT")
	err := endless.ListenAndServe(":"+appPort, r)
	if err != nil {
		fmt.Println("API Error:", err.Error())
	}
}
