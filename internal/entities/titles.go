package entities

type Title struct {
	Kind       string
	Content    string
	Chapter    string
	Last_verse string
	Footnote   []Footnote
	Crossref   []Crossref
}
