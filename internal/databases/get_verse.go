package databases

import (
	"context"
	"errors"
	"fmt"
	ent "github.com/BinaryGhost/verse-now/internal/entities"
	prs "github.com/BinaryGhost/verse-now/internal/parsers"
	"go.mongodb.org/mongo-driver/v2/bson"
	"sync"
)

// TODO: Keep verse-numbers like "1-2" in mind

func (db *Bible_db) ComposeVerses(ctx context.Context, acc *ent.WholeVerse, trans_abbr string, book_code string, alf *prs.AllReferences) error {
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

	var wg sync.WaitGroup
	err_chan := make(chan error, 10)

	for i < verse_range_len {
		wg.Add(1)
		cur := i

		go func(index int) {
			defer wg.Done()

			chapter_number := alf.Chapter

			lower_number := alf.MinVerses[index]
			higher_number := alf.MaxVerses[index]

			fmt.Println("lower_number: ", lower_number)
			fmt.Println("higher_number: ", higher_number)

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
						{Key: "verses.chapter", Value: chapter_number},
						// {Key: "verses.verse_min_range", Value: bson.D{
						// 	{Key: "$gte", Value: lower_number},
						// 	{Key: "$lte", Value: higher_number},
						// }}, // TODO: Handle verses, like 1-2
						{Key: "$and", Value: bson.A{
							bson.E{Key: "verses.verse_min_range", Value: bson.E{Key: "$gte", Value: lower_number}},
							bson.E{Key: "verses.verse_min_range", Value: bson.E{Key: "$lte", Value: higher_number}},
						},
						},
					}},
				},

				// {
				// 	{Key: "$match", Value: bson.D{
				// 		{Key: "$and", Value: bson.A{
				// 			bson.D{{Key: "verses.chapter", Value: chapter_number}},
				// 			bson.D{{Key: "verses.verse_min_range", Value: bson.D{
				// 				{Key: "$gte", Value: lower_number},
				// 				{Key: "$lte", Value: higher_number},
				// 			}}},
				// 		}},
				// 	}},
				// },
				{
					{Key: "$replaceRoot", Value: bson.D{
						{Key: "newRoot", Value: "$verses"},
					}},
				},
			}

			cursor, err := base_collection.Aggregate(ctx, pipeline)
			if err != nil {
				err_chan <- err
			}
			defer cursor.Close(ctx)

			var results []ent.Verse
			if err := cursor.All(ctx, &results); err != nil {
				err_chan <- err
			}

			acc.Verses = append(acc.Verses, results...)

		}(cur)

		//

		i++
	}

	wg.Wait()
	close(err_chan)

	for err := range err_chan {
		if err != nil {
			return err
		}
	}

	return nil
}
