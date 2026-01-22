package entities

type Footnote struct {
	References         string
	References_chapter uint64
	References_verse   Verse_number
	Text               string
}

type Crossref struct {
	References         string
	Belongs_to_chapter uint64
	Belongs_to_verse   Verse_number
	Text               string
}
