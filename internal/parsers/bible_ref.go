package parsers

import (
	"errors"
	"strings"
)

// Inside the url a string from a route with the following format
//
// Simpler thing
// <chapter>,<verses>
//
// Verse-Range
// <chapter>,<verses>-<verses>
//
// Verse-Range
// <chapter>,<verses>-<verses>.<verse>
//
// Verse-Range with another verse
// <chapter>,<verses>-<verses>
type BibleReferencesURL = string

// Chapter says, what it's name implies
//
// MinVerses/MaxVerses specify a range, they are adjacent to another:
// e.g: '1-2.3-4.5-6', would be
//
//	min [1, 3, 5]
//	max [2, 4, 6]
type AllReferences struct {
	Chapter   string
	MinVerses []string
	MaxVerses []string
}

// Validation does not happen here
func ReturnAllReferences(brfu BibleReferencesURL) (AllReferences, error) {
	_part := strings.Split(brfu, ",")
	chapter_part := _part[0]
	verse_part := _part[1]

	all_refs := AllReferences{Chapter: chapter_part}

	verses_dots := strings.Split(verse_part, ".")
	for _, vd := range verses_dots {
		vd_dash := strings.Split(vd, "-")

		if len(vd_dash) <= 1 {
			all_refs.MinVerses = append(all_refs.MinVerses, vd_dash[0])
			all_refs.MaxVerses = append(all_refs.MaxVerses, vd_dash[0])
		} else if len(vd_dash) == 2 {
			all_refs.MinVerses = append(all_refs.MinVerses, vd_dash[0])
			all_refs.MaxVerses = append(all_refs.MaxVerses, vd_dash[1])
		} else {
			return all_refs, errors.New("More than one '-' found")
		}
	}

	return all_refs, nil
}
