package entities

type Verse struct {
	global_locator         string `bson:"global_locator"`
	chapter                string `bson:"chapter"`
	verse_number           string `bson:"verse_number"`
	alternate_verse_number string `bson:""` // carefull
	text                   string `bson:"text"`
	is_a_list_element      bool   `bson:""` // carefull
	position               string `bson:"position"`
}
