package databases

import (
	"context"
	"errors"
	"fmt"
	ent "github.com/BinaryGhost/verse-now/internal/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// TODO: Keep verse-numbers like "1-2" in mind

func (db *Bible_db) ComposeVerses(ctx context.Context, acc *ent.WholeVerse, abbr string, book string, chapter string) error {
	base_collection := db.Collection(abbr)
	if base_collection == nil {
		error_str := fmt.Sprintf("Could not find collection of '%s'", abbr)
		return errors.New(error_str)
	}

	filter := bson.D{
		bson.E{Key: "general.about_book.bookname_in_english", Value: book},
		bson.E{Key: "verses", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$limit", Value: 1},
		},
		{
			{Key: "$unwind", Value: "$verses"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "verses.chapter", Value: chapter},
				{Key: "verses.", Value: chapter},
			}},
		},
		{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: "$verses"},
			}},
		},
	}

	cursor, err := base_collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []ent.Verse
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}
	acc.Verses = results

	return nil
}
