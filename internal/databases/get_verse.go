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

			// fmt.Println("chapter_number: ", chapter_number)
			// fmt.Println("lower_number: ", lower_number)
			// fmt.Println("higher_number: ", higher_number)

			pipeline := []bson.D{
				{
					{Key: "$match", Value: filter},
				},
				{
					{Key: "$limit", Value: 1},
				},
				{
					{Key: "$unwind", Value: "$content"},
				},
				{
					{Key: "$match", Value: bson.D{
						{Key: "$or", Value: bson.A{
							bson.D{{"content.role", "verse"}},
							bson.D{{"content.role", "footnote"}},
							bson.D{{"content.role", "crossref"}},

							// NOTE:
							// A table can contain verses, that are not directly requested,
							// but it does not make sense to exclude them, because you generally
							// want the context of an entire table, rather it being cut off strictly
							bson.D{{"content.role", "table"}},
						}},
						{Key: "content.chapter", Value: chapter_number},
						{Key: "$and", Value: bson.A{
							bson.D{{"content.verse_min_range", bson.D{{"$gte", lower_number}}}},
							bson.D{{"content.verse_min_range", bson.D{{"$lte", higher_number}}}},
						}},
						// TODO: Handle verses, like 1-2
					}},
				},
				{
					{Key: "$replaceRoot", Value: bson.D{
						{Key: "newRoot", Value: "$content"},
					}},
				},
			}

			cursor, err := base_collection.Aggregate(ctx, pipeline)
			if err != nil {
				err_chan <- err
			}
			defer cursor.Close(ctx)

			var results []any
			if err := cursor.All(ctx, &results); err != nil {
				err_chan <- err
			}

			acc.Anything = append(acc.Anything, results...)
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
