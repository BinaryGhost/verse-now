package databases

import (
	"context"
	"errors"
	"fmt"
	ent "github.com/BinaryGhost/verse-now/internal/entities"
	prs "github.com/BinaryGhost/verse-now/internal/parsers"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// TODO: Keep verse-numbers like "1-2" in mind

func (db *Bible_db) ComposeVerse(ctx context.Context, acc *ent.WholeVerse, trans_abbr string, book_code string, alf *prs.AllReferences) error {
	base_collection := db.Collection(trans_abbr)
	if base_collection == nil {
		error_str := fmt.Sprintf("Could not find collection of '%s'", trans_abbr)
		return errors.New(error_str)
	}

	filter := bson.D{
		bson.E{Key: "general.about_book.book_code", Value: book_code},
		bson.E{Key: "verses", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	verse_range_len := len(alf.MaxVerses)
	i := 0

	for i < verse_range_len {
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
					{Key: "verses.chapter", Value: alf.Chapter},
					{
						Key: "verses.verse_number",
						Value: bson.D{
							{Key: "$gte", Value: alf.MinVerses[i]},
							{Key: "$lte", Value: alf.MaxVerses[i]},
						},
					}, // TODO: Handle verses, like 1-2
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

		acc.Verses = append(acc.Verses, results...)

		//

		i++
	}

	return nil
}
