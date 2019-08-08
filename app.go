package main

import (
	"ecr-reunion/db"
	"ecr-reunion/facebook"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"os"
)

var appPort = os.Getenv("APP_PORT")
var ginIsReleaseMode = os.Getenv("APP_GIN_MODE_RELEASE") == "true"

func main() {
	if ginIsReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(db.ConnectMiddleware())

	// Add Facebook OAuth
	facebook.AuthController(r)

	err := endless.ListenAndServe(":"+appPort, r)
	if err != nil {
		fmt.Println("API Error:", err.Error())
	}
}
