package facebook

import (
	GinPassportFacebook "github.com/durango/gin-passport-facebook"
	"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"os"
)

var isDev = os.Getenv("APP_DEV_MODE") == "true"
var appPort = os.Getenv("APP_PORT")
var facebookClientID = os.Getenv("APP_FACEBOOK_CLIENT_ID")
var facebookClientSecret = os.Getenv("APP_FACEBOOK_CLIENT_SECRET")

func AuthController(router *gin.Engine) {
	redirectUrl := "https://api.classof09.org"
	if isDev {
		redirectUrl = "http://localhost:" + appPort
	}
	opts := &oauth2.Config{
		RedirectURL:  redirectUrl + "/oauth/callback",
		ClientID:     facebookClientID,
		ClientSecret: facebookClientSecret,
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}

	auth := router.Group("/auth/facebook")

	GinPassportFacebook.Routes(opts, auth)

	//auth.GET("/callback", GinPassportFacebook.Middleware(), func(c *gin.Context) {
	//	db := c.MustGet("db").(*mongo.Database)
	//	users := db.Collection("users")
	//
	//	user, err := GinPassportFacebook.GetProfile(c)
	//	if err != nil {
	//		log.Fatal(err.Error())
	//	}
	//
	//	u := models.User{
	//		Name:       user.Name,
	//		Email:      user.Email,
	//		FacebookId: user.Id,
	//	}
	//
	//	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	//	defer cancel()
	//
	//	_, err = users.InsertOne(ctx, u)
	//	if err != nil {
	//		log.Panic(err.Error())
	//	}
	//})
}
