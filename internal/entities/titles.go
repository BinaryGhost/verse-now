package entities

type Title struct {
	Kind       string
	Content    string
	Chapter    uint64
	Last_verse uint64
	Footnote   []Footnote
	Crossref   []Crossref
}
