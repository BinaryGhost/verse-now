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
	gather_about := []Gather{
		verses{}, footnotes{}, crossrefs{}, tables{}, titles{}, special_elements{},
	}
	var acc = ent.Chapter{}

	base_collection := db.Collection(trans_abbr)
	if base_collection == nil {
		error_str := fmt.Sprintf("Could not find collection of '%s'", trans_abbr)
		return errors.New(error_str)
	}

	if err := CollectAll(ctx, base_collection, gather_about, book_code, chapter, &acc); err != nil {
		return err
	}
	fmt.Println(len(acc.Verses))
	fmt.Println(len(acc.Crossrefs))
	fmt.Println(len(acc.Footnotes))
	fmt.Println(len(acc.Specials))
	fmt.Println(len(acc.Tables))
	fmt.Println(len(acc.Titles))

	return nil
}

func CollectAll(ctx context.Context, coll *mongo.Collection, gatherers []Gather, book_code string, chapter string, acc *ent.Chapter) error {
	var wg sync.WaitGroup
	err_chan := make(chan error, 8)

	for _, gatherer := range gatherers {
		wg.Add(1)
		go func(g Gather) {
			defer wg.Done()
			err := g.Gather(ctx, coll, book_code, chapter, acc)
			if err != nil {
				err_chan <- err
				return
			}
		}(gatherer)
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

type Gather interface {
	Gather(ctx context.Context, coll *mongo.Collection, book string, chapter string, acc *ent.Chapter) error
}

//
// NOTE: We need to set $limit to 1, since multiple (identical) documents can be fetched with this. I dont know why this happens
//

type verses struct{}

func (v verses) Gather(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, acc *ent.Chapter) error {
	filter := bson.D{
		bson.E{Key: "general.about_book.book_code", Value: book_code},
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
			}},
		},
		{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: "$verses"},
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
	acc.Verses = results

	return nil
}

type footnotes struct{}

func (f footnotes) Gather(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, acc *ent.Chapter) error {
	filter := bson.D{
		bson.E{Key: "general.about_book.book_code", Value: book_code},
		bson.E{Key: "footnotes", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$limit", Value: 1},
		},
		{
			{Key: "$unwind", Value: "$footnotes"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "footnotes.chapter", Value: chapter},
			}},
		},
		{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: "$footnotes"},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []ent.Footnote
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}
	acc.Footnotes = results

	return nil
}

type crossrefs struct{}

func (c crossrefs) Gather(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, acc *ent.Chapter) error {
	filter := bson.D{
		bson.E{Key: "general.about_book.book_code", Value: book_code},
		bson.E{Key: "cross_references", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$limit", Value: 1},
		},
		{
			{Key: "$unwind", Value: "$cross_references"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "cross_references.chapter", Value: chapter},
			}},
		},
		{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: "$cross_references"},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []ent.Crossref
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}
	acc.Crossrefs = results

	return nil
}

type titles struct{}

func (t titles) Gather(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, acc *ent.Chapter) error {
	filter := bson.D{
		bson.E{Key: "general.about_book.book_code", Value: book_code},
		bson.E{Key: "titles", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$limit", Value: 1},
		},
		{
			{Key: "$unwind", Value: "$titles"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "titles.chapter", Value: chapter},
			}},
		},
		{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: "$titles"},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []ent.Title
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}
	acc.Titles = results

	return nil
}

type tables struct{}

func (tb tables) Gather(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, acc *ent.Chapter) error {
	filter := bson.D{
		bson.E{Key: "general.about_book.book_code", Value: book_code},
		bson.E{Key: "tables", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$limit", Value: 1},
		},
		{
			{Key: "$unwind", Value: "$tables"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "tables.chapter", Value: chapter},
			}},
		},
		{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: "$tables"},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []ent.Table
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}
	acc.Tables = results

	return nil
}

type special_elements struct{}

func (s special_elements) Gather(ctx context.Context, coll *mongo.Collection, book_code string, chapter string, acc *ent.Chapter) error {
	filter := bson.D{
		bson.E{Key: "general.about_book.book_code", Value: book_code},
		bson.E{Key: "special_elems.specials", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$limit", Value: 1},
		},
		{
			{Key: "$unwind", Value: "$special"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "special.chapter", Value: chapter},
			}},
		},
		{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: "$special"},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []ent.Special
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}
	acc.Specials = results

	return nil
}
