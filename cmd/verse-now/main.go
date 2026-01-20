package main

import (
	"context"
	"fmt"

	db "github.com/BinaryGhost/verse-now/internal/databases"
	"github.com/gin-gonic/gin"
	// "net/http"
)

func main() {
	client := db.Client()
	bible_db := client.BibleDB()
	ctx, cancel := context.WithCancel(context.Background())

	defer client.Close(ctx)
	defer cancel()

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		// thing := c.Param("thing")
		// c.JSON(http.StatusOK, gin.H{
		// 	"message": thing,
		// })db.findOne
		fmt.Println(bible_db.ComposeChapter("Genesis", "1", ctx, "NLDNBG"))
	})

	// r.GET("/translations/:abbr/books/:book_name/chapters/:chapter_num")
	//
	// r.GET("/translations/:abbr/books/:book_name/chapters/:chapter_num/verses/:verses_notation")

	r.Run(":8080")
}
