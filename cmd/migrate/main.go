package main

import (
	"context"
	config "github.com/BinaryGhost/verse-now/configs"
	db "github.com/BinaryGhost/verse-now/internal/databases"
)

func main() {
	client := db.Client()
	db.IterateThroughBibleCollection(config.BABIJ_SOURCE, &client, context.Background())
}
