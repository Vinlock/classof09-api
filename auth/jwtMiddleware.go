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
	"log"
	"os"
	"time"
)

var identityKey = "_id"

func JWTMiddleware(r *gin.Engine) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "production",
		Key:        []byte(os.Getenv("APP_JWT_SECRET")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var facebookId string
			if err := c.ShouldBind(&facebookId); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			db := c.MustGet("database").(*mongo.Database)
			users := db.Collection("users")

			facebookUser, err := GinPassportFacebook.GetProfile(c)
			if err != nil {
				log.Fatal(err.Error())
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

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
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.User); ok && v.FacebookId != "" {
				return true
			}
			return false
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		LoginResponse:   nil,
		RefreshResponse: nil,
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				ID: claims[identityKey].(primitive.ObjectID),
			}
		},
		IdentityKey:   identityKey,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		CookieName:    "jwt",
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.GET(
		"/oauth/callback",
		GinPassportFacebook.Middleware(),
		authMiddleware.LoginHandler,
	)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/user", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			user, _ := c.Get(identityKey)
			c.JSON(200, gin.H{
				"userId": claims[identityKey],
				"name":   user.(*models.User).Name,
			})
		})
	}
}
