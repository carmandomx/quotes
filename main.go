package main

import (
	"github.com/carmandomx/quotes/controllers"
	"github.com/carmandomx/quotes/formatter"

	// "github.com/carmandomx/quotes/models"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	formatter.NewJSONFormatter()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/quotes", controllers.GetAllQuotes)
	r.PUT("/api/quotes/:id", controllers.UpdateQuote)
	r.POST("/api/quotes", controllers.CreateQuote)
	r.DELETE("/api/quotes/:id", controllers.DeleteQuote)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
