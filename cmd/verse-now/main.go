package main

import (
	"context"
	db "github.com/BinaryGhost/verse-now/internal/databases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	client := db.Client()
	defer client.Close(context.Background())

	r := gin.Default()
	r.GET("/echo/:thing", func(c *gin.Context) {
		thing := c.Param("thing")
		c.JSON(http.StatusOK, gin.H{
			"message": thing,
		})
	})

	r.GET("/translations/:abbr/books/:book_name/chapters/:chapter_num")

	r.GET("/translations/:abbr/books/:book_name/chapters/:chapter_num/verses/:verses_notation")

	r.Run(":8080")
}
