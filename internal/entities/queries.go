package entities

// TODO: Add mutexes

type Chapter struct {
	Verses    []Verse
	Footnotes []Footnote
	Crossrefs []Crossref
	Tables    []Table
	Specials  []Special
	Titles    []Title
}

type WholeVerse struct {
	Verses      []Verse
	InsideTable bool
	Footnotes   []Footnote
	Crossrefs   []Crossref
}
