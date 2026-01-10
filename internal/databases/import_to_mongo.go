package databases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TranslationStructure struct {
	General      any
	SpecialElems any
	Verses       any
	Titles       any
	Footnotes    any
	Tables       any
	CrossRefs    any
}

/*
Walk through a directory, where the babij is located ($PATH/babij_repo/collection) and find
babij-documents as $PATH/babij_repo/collection/.../*.json and use their translation_abbr as
collection name to be inserted inside of "bible_db".

Important: Do NOT insert multiple babij-collections, because collection can have translations,
with the same translation_abbr, thus making the collection corrupted.
*/
func IterateThroughBibleCollection(collection_path string, client *MClient) {
	db := client.BibleDB()

	err := filepath.Walk(collection_path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			exportBabijTranslation(path, db)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Something went wrong for '%s' :(\n", collection_path)
	}
}

func exportBabijTranslation(babij_translation_source string, bdb *Bible_db) {
	err := filepath.Walk(babij_translation_source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			var translation_abbr string
			base_name := strings.Split(filepath.Base(path), ".")[0]
			if strings.Contains(base_name, "not_possible") {
				fmt.Printf("Skipped '%s', because it has no content\n", path)
			}

			file, err := os.Open(path)
			if err != nil {
				log.Fatalf("Could not open '%s'", path)
			}

			read, _ := io.ReadAll(file)

			var book_json TranslationStructure
			if err := json.Unmarshal(read, &book_json); err != nil {
				log.Fatalf("JSON could not be parsed for '%s'", path)
			}

			if val := gjson.GetBytes(read, "general.about_translation.translation_abbr").String(); val == "" {
				log.Fatalf("Couldnt find translation_abbr for '%s'\n", path)
			} else {
				translation_abbr = val
			}

			// translation_abbr = book_json.General.AboutTranslation.TranslationAbbr

			coll := bdb.createCollection(translation_abbr)

			coll.insertBook(book_json, translation_abbr, base_name)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Could not go through '%s', because of error '%s'\n", babij_translation_source, err)
	}
}

type Bible_db struct {
	*mongo.Database
}

func (client *MClient) BibleDB() *Bible_db {
	return &Bible_db{client.mc.Database("bible_db")}
}

type translation_collection struct {
	*mongo.Collection
}

func (bdb *Bible_db) createCollection(trnl_abbr string) *translation_collection {
	return &translation_collection{bdb.Database.Collection(trnl_abbr)}
}

func (coll *translation_collection) insertBook(data any, trans_abbr string, base_name string) {
	str := fmt.Sprintf("collection{'%s'} -> book{%s}\n", trans_abbr, base_name)
	_, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatalf("DID NOT insert %s, because of %s\n", str, err)

	}
	fmt.Printf("Inserted into %s", str)
}
