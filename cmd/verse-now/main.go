package main

import (
	"context"
	"fmt"
	db "github.com/BinaryGhost/verse-now/internal/databases"
	ent "github.com/BinaryGhost/verse-now/internal/entities"
	prs "github.com/BinaryGhost/verse-now/internal/parsers"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	client := db.Client()
	bible_db := client.BibleDB()
	ctx, cancel := context.WithCancel(context.Background())

	defer client.Close(ctx)
	defer cancel()

	r := gin.Default()
	r.GET("/translations/:translation/books/:book/chapters/:chapter", func(c *gin.Context) {
		translation := c.Param("translation")
		book := c.Param("book")

		chapter, err := strconv.ParseUint(c.Param("chapter"), 10, 64)
		if err != nil || chapter == 0 {
			err_msg := fmt.Sprintf("Can not use '%d' as chapter-number", chapter)
			c.JSON(400, gin.H{
				"error": err_msg,
			})
		}

		todo := ent.Chapter{}

		bible_db.ComposeChapter(ctx, &todo, translation, book, chapter)
		fmt.Println(">>>>", len(todo.Anything))
	})

	r.GET("/translations/:translation/books/:book/reference/:reference", func(c *gin.Context) {
		translation := c.Param("translation")
		book := c.Param("book")
		bible_reference_string := prs.BibleReferencesURL(c.Param("reference"))

		make_reference, _ := prs.ReturnAllReferences(bible_reference_string)

		todo := ent.WholeVerse{}

		// bible_db.ComposeVerses(ctx, &todo, "NLDNBG", "PSA", &make_reference)
		bible_db.ComposeVerses(ctx, &todo, translation, book, &make_reference)

		fmt.Println(">>>>", len(todo.Anything))
	})

	// r.GET("/translations/:abbr/books/:book_name/chapters/:chapter_num")
	//
	// r.GET("/translations/:abbr/books/:book_name/chapters/:chapter_num/verses/:verses_notation")

	r.Run(":8080")
}
