package entities

type Title struct {
	Kind               string
	Content            string
	Chapter            uint64
	Verse_min_range    uint64
	Verse_max_range    uint64
	Verse_min_notation string
	Verse_max_notation string
	Footnotes          []Footnote
	Crossrefs          []Crossref
}
