package auth

import (
	"context"
	"ecr-reunion/models"
	"ecr-reunion/typeform"
	jwt "github.com/appleboy/gin-jwt"
	GinPassportFacebook "github.com/durango/gin-passport-facebook"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"log"
	"os"
	"time"
)

var identityKey = "_id"

func logout(c *gin.Context) {
	cookieDomain := os.Getenv("APP_COOKIE_DOMAIN")
	secureCookie := os.Getenv("APP_DEV_MODE") != "true"
	c.SetCookie("token", "", 0, "/", cookieDomain, secureCookie, false)
	redirectTo := c.Request.URL.Query().Get("redirect")
	c.Redirect(302, redirectTo)
}

func JWTMiddleware(r *gin.Engine) {
	authMiddleware := getMiddleware()

	authGroup := r.Group("/auth")
	GinPassportFacebook.Routes(getOauth2Config(), authGroup)
	authGroup.GET(
		"/callback",
		GinPassportFacebook.Middleware(),
		authMiddleware.LoginHandler,
	)
	authGroup.GET("/logout", logout)
	authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
	authGroup.GET("/user", authMiddleware.MiddlewareFunc(), getUser)
}

func getMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "production",
		Key:             []byte(os.Getenv("APP_JWT_SECRET")),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour * 24,
		Authenticator:   authenticator,
		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authorizator:    authorizator,
		Unauthorized:    unauthorized,
		LoginResponse:   loginResponse,
		RefreshResponse: nil,
		TokenLookup:     "header: Authorization",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
		CookieName:      "token",
		CookieDomain:    os.Getenv("APP_COOKIE_DOMAIN"),
		CookieHTTPOnly:  false,
		SendCookie:      true,
		SecureCookie:    os.Getenv("APP_DEV_MODE") != "true",
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}

func authenticator(c *gin.Context) (interface{}, error) {
	db := c.MustGet("database").(*mongo.Database)
	users := db.Collection("users")

	facebookUser, err := GinPassportFacebook.GetProfile(c)
	if err != nil {
		return "", jwt.ErrFailedAuthentication
	}

	ctx := context.Background()

	filter := bson.M{"facebookId": facebookUser.Id}
	var user models.User
	err = users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		user = models.User{
			Name:       facebookUser.Name,
			Email:      facebookUser.Email,
			FacebookId: facebookUser.Id,
		}
		_, err = users.InsertOne(ctx, user)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}
	}
	return user, nil
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(models.User); ok {
		return jwt.MapClaims{
			jwt.IdentityKey: v.FacebookId,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	db := c.MustGet("database").(*mongo.Database)
	users := db.Collection("users")
	ctx := context.Background()

	claims := jwt.ExtractClaims(c)

	facebookId := claims[jwt.IdentityKey].(string)

	filter := bson.M{"facebookId": facebookId}
	var user models.User
	err := users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil
	}
	return user
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(models.User); ok && v.FacebookId != "" {
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":  code,
		"error": message,
	})
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.Redirect(302, os.Getenv("APP_FACEBOOK_LOGIN_REDIRECT"))
}

func getOauth2Config() *oauth2.Config {
	redirectUrl := "https://api.classof09.org"
	if os.Getenv("APP_DEV_MODE") == "true" {
		redirectUrl = "http://localhost:" + os.Getenv("APP_PORT")
	}

	opts := &oauth2.Config{
		RedirectURL:  redirectUrl + "/auth/callback",
		ClientID:     os.Getenv("APP_FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("APP_FACEBOOK_CLIENT_SECRET"),
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}

	return opts
}

func getUser(c *gin.Context) {
	db := c.MustGet("database").(*mongo.Database)
	users := db.Collection("users")
	ctx := context.Background()

	claims := jwt.ExtractClaims(c)

	facebookId := claims[jwt.IdentityKey].(string)

	filter := bson.M{"facebookId": facebookId}
	var user models.User
	err := users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	surveyDone := false
	typeformApi := typeform.NewTypeformApi(os.Getenv("APP_TYPEFORM_TOKEN"))
	params := typeform.GetResponsesParams{
		FormId: os.Getenv("APP_SURVEY1_ID"),
		Query:  facebookId,
	}
	response, err := typeformApi.GetResponses(params)
	if err != nil {
		surveyDone = false
	} else if response.TotalItems > 0 {
		surveyDone = true
	}

	if !surveyDone {
		c.SetCookie(
			"typeform_done",
			"",
			0,
			"/",
			os.Getenv("APP_COOKIE_DOMAIN"),
			os.Getenv("APP_DEV_MODE") != "true",
			false,
		)
	}

	// Get Total Entries
	totalEntries := 0
	onlyCompletedEntries := true
	params = typeform.GetResponsesParams{
		FormId:    "CWy6cX",
		Completed: &onlyCompletedEntries,
	}
	response, err = typeformApi.GetResponses(params)
	if err != nil {
		totalEntries = 0
	} else {
		totalEntries = response.TotalItems
	}

	c.JSON(200, gin.H{
		"id":           user.ID,
		"fbId":         user.FacebookId,
		"name":         user.Name,
		"surveyDone":   surveyDone,
		"totalEntries": totalEntries,
	})
}
