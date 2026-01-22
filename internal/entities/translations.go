package entities

type TranslationStructure struct {
	General          any
	Special_Elems    []Special // Note to myself: keys have to match (in letter, not case)
	Verses           []Verse
	Titles           []Title
	Footnotes        []Footnote
	Tables           []Table
	Cross_references []Crossref
}
