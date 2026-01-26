package databases

import (
	"context"
	"errors"
	"fmt"
	ent "github.com/BinaryGhost/verse-now/internal/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"sync"
)

// TODO: Keep verse-numbers like "1-2" in mind

func (db *Bible_db) ComposeChapter(book_code string, chapter string, ctx context.Context, trans_abbr string) error {
	var acc = ent.Chapter{}

	base_collection := db.Collection(trans_abbr)
	if base_collection == nil {
		error_str := fmt.Sprintf("Could not find collection of '%s'", trans_abbr)
		return errors.New(error_str)
	}

	if err := CollectAll(ctx, base_collection, book_code, chapter, &acc); err != nil {
		return err
	}
	// fmt.Println(len(acc.Verses))
	// fmt.Println(len(acc.Crossrefs))
	// fmt.Println(len(acc.Footnotes))
	// fmt.Println(len(acc.Specials))
	// fmt.Println(len(acc.Tables))
	// fmt.Println(len(acc.Titles))

	return nil
}

func CollectAll(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, acc *ent.Chapter) error {
	var wg sync.WaitGroup
	err_chan := make(chan error, 8)

	for _, role := range ent.ExistingRoles {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			err := gather(ctx, coll, book_code, chapter, role, acc)
			if err != nil {
				err_chan <- err
				return
			}
		}(role)
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

func gather(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, role string, acc *ent.Chapter) error {
	filter := bson.D{
		bson.E{Key: "general.about_book.book_code", Value: book_code},
	}

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
				{Key: "content.chapter", Value: chapter},
				{Key: "role", Value: role},
			}},
		},
		{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: "$content"},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []ent.Verse
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}
	acc.Anything = append(acc.Anything, results)

	return nil

}

//
// NOTE: We need to set $limit to 1, since multiple (identical) documents can be fetched with this. I dont know why this happens
//

// type verses struct{}
//
// func (v verses) Gather(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, acc *ent.Chapter) error {
// 	filter := bson.D{
// 		bson.E{Key: "general.about_book.book_code", Value: book_code},
// 		bson.E{Key: "verses", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
// 	}
//
// 	pipeline := []bson.D{
// 		{
// 			{Key: "$match", Value: filter},
// 		},
// 		{
// 			{Key: "$limit", Value: 1},
// 		},
// 		{
// 			{Key: "$unwind", Value: "$verses"},
// 		},
// 		{
// 			{Key: "$match", Value: bson.D{
// 				{Key: "verses.chapter", Value: chapter},
// 			}},
// 		},
// 		{
// 			{Key: "$replaceRoot", Value: bson.D{
// 				{Key: "newRoot", Value: "$verses"},
// 			}},
// 		},
// 	}
//
// 	cursor, err := coll.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return err
// 	}
// 	defer cursor.Close(ctx)
//
// 	var results []ent.Verse
// 	if err := cursor.All(ctx, &results); err != nil {
// 		return err
// 	}
// 	acc.Verses = results
//
// 	return nil
// }
