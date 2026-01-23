package entities

type Title struct {
	Kind       string
	Content    string
	Chapter    uint64
	Last_verse Vrs_number_strct
	Footnote   []Footnote
	Crossref   []Crossref
}
