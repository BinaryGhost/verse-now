package databases

import (
	"context"
	"errors"
	"fmt"
	ent "github.com/BinaryGhost/verse-now/internal/entities"
	"sync"

	// "github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	// "go.mongodb.org/mongo-driver/v2/mongo/options"
	// "go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (db *Bible_db) ComposeChapter(book string, chapter string, ctx context.Context, abbr string) error {
	// title
	// tables
	// special_elements

	gather_about := []Gather{
		verse{}, footnote{}, crossrefs{}, tables{}, titles{},
	}

	base_collection := db.Collection(abbr)
	if base_collection == nil {
		error_str := fmt.Sprintf("Could not find collection of '%s'", abbr)
		return errors.New(error_str)
	}

	_ = Collect(ctx, base_collection, gather_about, book, chapter)

	// _ = GatherVerses(ctx, base_collection, book, chapter)

	return nil
}

func Collect(ctx context.Context, coll *mongo.Collection, gatherers []Gather, book string, chapter string) error {
	var wg sync.WaitGroup
	err_chan := make(chan error, 8)

	for _, gatherer := range gatherers {
		wg.Add(1)
		go func(g Gather) {
			defer wg.Done()
			_, err := g.Gather(ctx, coll, book, chapter)
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
	Gather(ctx context.Context, coll *mongo.Collection, book string, chapter string) (any, error)
}

type verse struct{}

func (v verse) Gather(ctx context.Context, coll *mongo.Collection, book string, chapter string) (any, error) {
	filter := bson.D{
		bson.E{Key: "general.about_book.bookname_in_english", Value: book},
		bson.E{Key: "verses", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
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
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ent.Verse
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for _, verse := range results {
		fmt.Printf("Text: %s, Chapter: %s, VerseNumber: %s\n", verse.Text, verse.Chapter, verse.Verse_number)
	}

	return results, nil
}

type footnote struct{}

func (f footnote) Gather(ctx context.Context, coll *mongo.Collection, book string, chapter string) (any, error) {
	filter := bson.D{
		bson.E{Key: "general.about_book.bookname_in_english", Value: book},
		bson.E{Key: "footnotes", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$unwind", Value: "$footnotes"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "footnotes.references_chapter", Value: chapter},
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
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ent.Footnote
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for _, footnote := range results {
		fmt.Printf("FOOTNOTE: %s, fChapter: %s, fNumber: %s\n", footnote.Text, footnote.References_chapter, footnote.References)
	}

	return results, nil
}

type crossrefs struct{}

func (c crossrefs) Gather(ctx context.Context, coll *mongo.Collection, book string, chapter string) (any, error) {
	filter := bson.D{
		bson.E{Key: "general.about_book.bookname_in_english", Value: book},
		bson.E{Key: "cross_references", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$unwind", Value: "$cross_references"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "cross_references.belongs_to_chapter", Value: chapter},
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
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ent.Crossref
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for _, crossref := range results {
		fmt.Printf("CROSSREFERENCE: %s, cChapter: %s, cNumber: %s\n", crossref.Text, crossref.Belongs_to_chapter, crossref.References)
	}

	return results, nil
}

type titles struct{}

func (t titles) Gather(ctx context.Context, coll *mongo.Collection, book string, chapter string) (any, error) {
	filter := bson.D{
		bson.E{Key: "general.about_book.bookname_in_english", Value: book},
		bson.E{Key: "titles", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
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
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ent.Title
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for _, titles := range results {
		fmt.Printf("TITLES: %s, tChapter: %s, tNumber: %s\n", titles.Content, titles.Chapter, titles.Last_verse)
	}

	return results, nil
}

type tables struct{}

func (tb tables) Gather(ctx context.Context, coll *mongo.Collection, book string, chapter string) (any, error) {
	filter := bson.D{
		bson.E{Key: "general.about_book.bookname_in_english", Value: book},
		bson.E{Key: "tables", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
	}

	pipeline := []bson.D{
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$unwind", Value: "$tables"},
		},
		{
			{Key: "$match", Value: bson.D{
				{Key: "tables.last_chapter", Value: chapter},
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
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ent.Table
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for _, tables := range results {
		fmt.Println(tables.String())
	}

	return results, nil
}
