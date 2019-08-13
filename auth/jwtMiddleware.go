package auth

import (
	"context"
	"ecr-reunion/models"
	jwt "github.com/appleboy/gin-jwt"
	GinPassportFacebook "github.com/durango/gin-passport-facebook"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"log"
	"os"
	"time"
)

var identityKey = "_id"

func JWTMiddleware(r *gin.Engine) {
	authMiddleware := getMiddleware()

	authGroup := r.Group("/auth")

	GinPassportFacebook.Routes(getOauth2Config(), authGroup)

	authGroup.GET(
		"/callback",
		GinPassportFacebook.Middleware(),
		authMiddleware.LoginHandler,
	)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{
			"code":  "PAGE_NOT_FOUND",
			"error": "Page not found",
		})
	})

	authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
	authGroup.GET("/user", authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		db := c.MustGet("database").(*mongo.Database)
		users := db.Collection("users")
		ctx := context.Background()

		claims := jwt.ExtractClaims(c)

		claimId := claims[jwt.IdentityKey].(string)
		id, err := primitive.ObjectIDFromHex(claimId)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}

		filter := bson.M{"_id": id}
		var user models.User
		err = users.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}

		responses := db.Collection("responses")
		responseFilter := bson.M{"form_response.hidden.id": id}
		var response interface{}
		err = responses.FindOne(ctx, responseFilter).Decode(response)
		surveyDone := true
		if err != nil {
			surveyDone = false
		}
		c.JSON(200, gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"surveyDone": surveyDone,
		})
	})
}

func getMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "production",
		Key:        []byte(os.Getenv("APP_JWT_SECRET")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(c *gin.Context) (interface{}, error) {
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
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(models.User); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			db := c.MustGet("database").(*mongo.Database)
			users := db.Collection("users")
			ctx := context.Background()

			claims := jwt.ExtractClaims(c)

			claimId := claims[jwt.IdentityKey].(string)
			id, err := primitive.ObjectIDFromHex(claimId)
			if err != nil {
				return nil
			}

			filter := bson.M{"_id": id}
			var user models.User
			err = users.FindOne(ctx, filter).Decode(&user)
			if err != nil {
				return nil
			}
			return user
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(models.User); ok && v.FacebookId != "" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":  code,
				"error": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.Redirect(302, "http://localhost:8000")
		},
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
