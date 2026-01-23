package entities

type Footnote struct {
	References         string
	References_chapter uint64
	References_verse   Vrs_number_strct
	Text               string
}

type Crossref struct {
	References         string
	Belongs_to_chapter uint64
	Belongs_to_verse   Vrs_number_strct
	Text               string
}
