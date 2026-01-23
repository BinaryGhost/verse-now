package entities

type Verse struct {
	Global_locator         string
	Chapter                uint64
	Alternate_verse_number string // It highlights, that the verse-number can also be different. However, it does not effect anything, and should not be used if empty
	Text                   string
	Is_a_list_element      bool
	Position               string
	Verse_min_range        uint64
	Verse_max_range        uint64
	Verse_min_notation     string
	Verse_max_notation     string
}
