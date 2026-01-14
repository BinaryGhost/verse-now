package entities

type Verse struct {
	Global_locator         string
	Chapter                string
	Verse_number           string
	Alternate_verse_number string // carefull
	Text                   string
	Is_a_list_element      bool // carefull
	Position               string
}
