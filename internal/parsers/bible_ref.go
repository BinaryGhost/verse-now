package parsers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Inside the url a string from a route with the following format
//
// Simpler thing
// <chapter>:<verses>
//
// Verse-Range
// <chapter>:<verses>-<verses>
//
// Verse-Range
// <chapter>:<verses>-<verses>.<verse>
//
// Verse-Range with another verse
// <chapter>:<verses>-<verses>
type BibleReferencesURL = string

// Chapter says, what it's name implies
//
// MinVerses/MaxVerses specify a range, they are adjacent to another:
// e.g: '1-2.3-4.5-6', would be
//
//	min [1, 3, 5]
//	max [2, 4, 6]
type AllReferences struct {
	Chapter   uint64
	MinVerses []uint64
	MaxVerses []uint64
}

func ReturnAllReferences(brfu BibleReferencesURL) (AllReferences, error) {
	_part := strings.Split(brfu, ":")
	chapter_part := _part[0]
	verse_part := _part[1]

	chapter, err := strconv.ParseUint(chapter_part, 10, 64)
	if err != nil {
		return AllReferences{}, errors.New("Invalid chapter-value")
	}

	all_refs := AllReferences{Chapter: chapter}

	verses_dots := strings.Split(verse_part, ".")
	for _, vd := range verses_dots {
		vd_dash := strings.Split(vd, "-")

		if len(vd_dash) > 2 {
			return all_refs, errors.New("Invalid amound of '-'")
		} else {
			min, err := strconv.ParseUint(vd_dash[0], 10, 64)
			if err != nil {
				return all_refs, errors.New("Could not make starting verse into a number")
			}

			var max uint64
			if len(vd_dash) == 1 { // exactly one verse
				max = min
			} else {
				max, err = strconv.ParseUint(vd_dash[1], 10, 64)
				if err != nil {
					return all_refs, errors.New("Could not make starting verse into a number")
				}
			}

			if min == 0 || max == 0 {
				return all_refs, errors.New("Verse is zero")
			}

			if min > max {
				str := fmt.Sprintf("Invalid reference found '%d-%d'", min, max)
				return all_refs, errors.New(str)
			}

			all_refs.MinVerses = append(all_refs.MinVerses, min)
			all_refs.MaxVerses = append(all_refs.MaxVerses, max)
		}
	}

	return all_refs, nil
}
