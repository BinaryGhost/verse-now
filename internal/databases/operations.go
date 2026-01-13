package databases

import (
	"context"
	"errors"
	"fmt"

	// "github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	// "go.mongodb.org/mongo-driver/v2/mongo/options"
	// "go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (db *Bible_db) ComposeChapter(book string, chapter string, ctx context.Context, abbr string) error {
	// title
	// footnotes
	// tables
	// special_elements
	// verses
	// crossrefs

	base_collection := db.Collection(abbr)
	if base_collection == nil {
		error_str := fmt.Sprintf("Could not find collection of '%s'", abbr)
		return errors.New(error_str)
	}

	inside_book := bson.E{Key: "general.about_book.bookname_in_english", Value: book}

	verses, _ := Gather(
		ctx,
		base_collection,
		Filter{
			filtering_opts: bson.D{
				inside_book,
				{Key: "verses", Value: bson.D{{Key: "$ne", Value: bson.A{}}}},
			},
			bookname:   book,
			query_kind: "verses",
		},
	)
	fmt.Println(verses)

	return nil
}

type Filter struct {
	filtering_opts bson.D
	bookname       string
	query_kind     string
}

func Gather(ctx context.Context, coll *mongo.Collection, filter Filter) (bson.M, error) {
	cursor := coll.FindOne(ctx, filter.filtering_opts)

	if cursor == nil {
		error_str := fmt.Sprintf("Bookname not found for '%s' of kind '%s'", filter.bookname, filter.query_kind)
		return nil, errors.New(error_str)
	}

	var result bson.M
	if err := cursor.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
