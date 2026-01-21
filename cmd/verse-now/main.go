package main

import (
	"context"
	db "github.com/BinaryGhost/verse-now/internal/databases"
	ent "github.com/BinaryGhost/verse-now/internal/entities"
	prs "github.com/BinaryGhost/verse-now/internal/parsers"
	"github.com/gin-gonic/gin"
)

func main() {
	client := db.Client()
	bible_db := client.BibleDB()
	ctx, cancel := context.WithCancel(context.Background())

	defer client.Close(ctx)
	defer cancel()

	r := gin.Default()
	r.GET("/test/:reference", func(c *gin.Context) {
		bible_reference_string := prs.BibleReferencesURL(c.Param("reference"))
		make_reference, _ := prs.ReturnAllReferences(bible_reference_string)

		todo := ent.WholeVerse{}

		bible_db.ComposeVerse(ctx, &todo, "NLDNBG", "PSA", &make_reference)
	})

	// r.GET("/translations/:abbr/books/:book_name/chapters/:chapter_num")
	//
	// r.GET("/translations/:abbr/books/:book_name/chapters/:chapter_num/verses/:verses_notation")

	r.Run(":8080")
}
