package main

import (
	"ecr-reunion/auth"
	"ecr-reunion/controllers"
	"ecr-reunion/db"
	"fmt"
	"github.com/bugsnag/bugsnag-go"
	bugsnagGin "github.com/bugsnag/bugsnag-go/gin"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

func main() {
	ginIsReleaseMode := os.Getenv("APP_GIN_MODE_RELEASE") == "true"
	if ginIsReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	bugsnag.Configure(bugsnag.Configuration{
		APIKey: os.Getenv("APP_BUGSNAG_API_KEY"),
		// The import paths for the Go packages containing your source files
		ProjectPackages: []string{"main", "ecr-reunion"},
		AppVersion:      os.Getenv("APP_COMMIT_SHA"),
	})

	r := gin.Default()

	r.Use(bugsnagGin.AutoNotify())

	r.GET("/test", func(c *gin.Context) {
		err := bugsnag.Notify(fmt.Errorf("Test error"))
		if err != nil {
			log.Panic("BUGSNAG_ERROR")
		}
		c.JSON(200, gin.H{})
	})

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
	controllers.TypeformController(r)

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
