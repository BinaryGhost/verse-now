package entities

type Footnote struct {
	References         string
	References_chapter string
	References_verse   string
	Text               string
}

type Crossref struct {
	References         string
	Belongs_to_chapter string
	Belongs_to_verse   string
	Text               string
}
