package entities

type Footnote struct {
	references         string `bson:"references"`
	references_chapter string `bson:"references_chapter"`
	references_verse   string `bson:"references_verse"`
	text               string `bson:"text"`
}

type Crossref struct {
	references         string `bson:"references"`
	belongs_to_chapter string `bson:"belongs_to_chapter"`
	belongs_to_verse   string `bson:"belongs_to_verse"`
	text               string `bson:"text"`
}
