package server

import (
	"log"
	"os"
	"time"

	models "github.com/alixleger/open-flight-core/db"
	"github.com/alixleger/open-flight-core/server/handlers"
	"github.com/alixleger/open-flight-core/services/skyscanner"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/influxdata/influxdb/client/v2"
)

// Server type
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

// New function is the server constructor
func New(db *gorm.DB, skyscannerClient *skyscanner.Client, influxdbClient *client.Client) *Server {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("influxdbClient", influxdbClient)
		c.Set("skyscannerClient", skyscannerClient)
		c.Next()
	})

	authMiddleware, err := getAuthMiddleware()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.GET("/flights", handlers.GetFlights)
	router.GET("/places", handlers.GetPlaces)
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
		auth.PATCH("/user", handlers.PatchUser)
		auth.GET("/flightprices", handlers.GetFlightPrices)
		auth.GET("/favflights", handlers.GetFavFlights)
		auth.POST("/favflight", handlers.PostFavFlight)
	}

	return &Server{DB: db, Router: router}
}

func getAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "openflight auth zone",
		Key:         []byte(os.Getenv("SECRET_KEY")),
		Timeout:     time.Hour * time.Duration(24),
		MaxRefresh:  time.Hour * time.Duration(24*7),
		IdentityKey: handlers.IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					handlers.IdentityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			db := c.MustGet("db").(*gorm.DB)
			var user models.User
			db.Where("email = ?", claims[handlers.IdentityKey]).First(&user)
			return &user
		},
		Authenticator: handlers.Authenticate,
		TimeFunc:      time.Now,
	})
}

// Run api server
func (server *Server) Run() {
	server.Router.Run(":" + os.Getenv("PORT"))
}
