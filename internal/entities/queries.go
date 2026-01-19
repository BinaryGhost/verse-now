package entities

type Chapter struct {
	Verses    []Verse
	Footnotes []Footnote
	Crossrefs []Crossref
	Tables    []Table
	Specials  []Special
	Titles    []Title
}
