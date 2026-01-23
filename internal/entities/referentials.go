package entities

type Footnote struct {
	References         string
	Chapter            uint64
	Verse_min_range    uint64
	Verse_max_range    uint64
	Verse_min_notation string
	Verse_max_notation string
	Text               string
}

type Crossref struct {
	References         string
	Chapter            uint64
	Verse_min_range    uint64
	Verse_max_range    uint64
	Verse_min_notation string
	Verse_max_notation string
	Text               string
}
