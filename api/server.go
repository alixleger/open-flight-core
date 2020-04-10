package api

import (
	"log"
	"os"
	"time"

	"github.com/alixleger/open-flight/back/api/handlers"
	"github.com/alixleger/open-flight/back/api/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var identityKey = "id"

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims[identityKey],
		"text":   "Hello World.",
	})
}

// Server type
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

// New function is the server constructor
func New(db *gorm.DB) Server {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Router use DB
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	authMiddleware, err := getAuthMiddleware()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.POST("/register", handlers.Register)
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/logout", authMiddleware.LogoutHandler)

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "NOT_FOUND", "message": "Ressource not found"})
	})

	auth := router.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// TODO: List routes which need authentication
		auth.GET("/hello", helloHandler)
	}

	return Server{DB: db, Router: router}
}

func getAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "openflight auth zone",
		Key:         []byte(os.Getenv("SECRET_KEY")),
		Timeout:     time.Hour * time.Duration(24),
		MaxRefresh:  time.Hour * time.Duration(24*7),
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: handlers.Authenticate,
		TimeFunc:      time.Now,
	})
}

// Run api server
func (server *Server) Run() {
	server.Router.Run(":" + os.Getenv("PORT"))
}
