package entities

// Min-Max specify a verse like '1-2', but can also be just '1', if it isnt a range
//
// For some translations '1a' or '1b' exist, for those, we just set their value accordingly
// 1b => {min: 1, max: 1, min: "a", "b"}
type Vrs_number_strct struct {
	min_range    uint64
	max_range    uint64
	min_notation string
	max_notation string
}

type Verse struct {
	Global_locator         string
	Chapter                uint64
	Verse_number           Vrs_number_strct
	Alternate_verse_number string // carefull
	Text                   string
	Is_a_list_element      bool // carefull
	Position               string
}
