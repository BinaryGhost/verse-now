package entities

type Special struct {
	Role               string // special
	Kind               string
	Content            string
	Chapter            uint64
	Verse_min_range    uint64
	Verse_max_range    uint64
	Verse_min_notation string
	Verse_max_notation string
}
