package api

import (
	"github.com/alixleger/open-flight/back/api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Server type
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

// New function is the server constructor
func New(db *gorm.DB) Server {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	router.GET("/books", controllers.FindBooks)
	router.GET("/books/:id", controllers.FindBook)
	router.POST("/books", controllers.CreateBook)
	router.PATCH("/books/:id", controllers.UpdateBook)
	router.DELETE("/books/:id", controllers.DeleteBook)

	return Server{DB: db, Router: router}
}

// Run api server
func (server *Server) Run() {
	server.Router.Run(":8000")
}
