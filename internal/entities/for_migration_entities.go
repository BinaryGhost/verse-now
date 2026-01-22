package entities

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func convert_num(num string) (uint64, error) {
	uinti, err := strconv.ParseUint(num, 10, 64)

	// Some things like footnotes can appear the beginning of chapter,
	// thus no 0 evalutation
	if err != nil {
		err_msg := fmt.Sprintf("num is 0 or is '%s'!", err)
		return 0, errors.New(err_msg)
	}

	return uinti, nil
}

func convert_to_right_footnote() {
}

func convert_to_right_crossref() {
}

func convert_to_verse_number(verse string) (Verse_number, error) {
	if strings.Contains(verse, "-") {
		_parts := strings.Split(verse, "-")

		left, err := strconv.ParseUint(_parts[0], 10, 64)
		if err != nil {
			return Verse_number{}, errors.New("Left 'number' not number")
		}

		right, err := strconv.ParseUint(_parts[1], 10, 64)
		if err != nil {
			return Verse_number{}, errors.New("Right 'number' not number")
		}

		return Verse_number{
			min: left,
			max: right,
		}, nil
	}

	re := regexp.MustCompile(`^(\d+)([a-zA-Z]+)$`)
	matches := re.FindStringSubmatch(verse)

	if len(matches) == 3 {
		num, err := convert_num(matches[1])
		if err != nil {
			return Verse_number{}, errors.New("Notation number wrong")
		}
		return Verse_number{
			min:          num,
			max:          num,
			min_notation: matches[2],
			max_notation: matches[2],
		}, nil
	}

	num, err := convert_num(verse)
	if err != nil {
		return Verse_number{}, err
	}

	return Verse_number{
		min: num,
		max: num,
	}, nil
}

func Migration_to_real_Structure(ms *Migration_Structure) (TranslationStructure, error) {
	tsst := TranslationStructure{}
	tsst.General = ms.General

	for _, v := range ms.Verses {
		chapt, err1 := convert_num(v.Chapter)
		vrs, err2 := convert_to_verse_number(v.Verse_number)

		if err1 != nil {
			err_mgs := fmt.Sprintf("%v, OH OH ~Verse", err1)
			return tsst, errors.New(err_mgs)
		}

		if err2 != nil {
			err_mgs := fmt.Sprintf("%v, OH OH not Verse_number again ~Verse", err2)
			return tsst, errors.New(err_mgs)
		}

		verse := Verse{
			Global_locator:         v.Global_locator,
			Chapter:                chapt,
			Verse_number:           vrs,
			Alternate_verse_number: v.Alternate_verse_number,
			Text:                   v.Text,
			Is_a_list_element:      v.Is_a_list_element,
			Position:               v.Position,
		}

		tsst.Verses = append(tsst.Verses, verse)
	}

	for _, f := range ms.Footnotes {
		chapt, err1 := convert_num(f.References_chapter)
		vrs, err2 := convert_to_verse_number(f.References_verse)

		if err1 != nil || err2 != nil {
			err_mgs := fmt.Sprintf("%s, %s, OH OH ~Footnotes", err1, err2)
			return tsst, errors.New(err_mgs)
		}

		footnote := Footnote{
			References:         f.References,
			References_chapter: chapt,
			References_verse:   vrs,
			Text:               f.Text,
		}

		tsst.Footnotes = append(tsst.Footnotes, footnote)
	}
	for _, cr := range ms.Crossrefs {
		chapt, err1 := convert_num(cr.Belongs_to_chapter)
		vrs, err2 := convert_to_verse_number(cr.Belongs_to_verse)

		if err1 != nil || err2 != nil {
			return tsst, errors.New("OH OH ~Crossrefs")
		}

		crossref := Crossref{
			References:         cr.References,
			Belongs_to_chapter: chapt,
			Belongs_to_verse:   vrs,
			Text:               cr.Text,
		}

		tsst.Cross_references = append(tsst.Cross_references, crossref)
	}

	for _, s := range ms.Raw_Specials.Specials {
		chapt, err1 := convert_num(s.Chapter)
		vrs, err2 := convert_num(s.Last_verse)

		if err1 != nil || err2 != nil {
			return tsst, errors.New("OH OH ~Specials")
		}

		special := Special{
			Kind:       s.Kind,
			Content:    s.Content,
			Chapter:    chapt,
			Last_verse: vrs,
		}

		tsst.Special_Elems = append(tsst.Special_Elems, special)
	}

	for _, t := range ms.Titles {
		chapt, err1 := convert_num(t.Chapter)
		vrs, err2 := convert_num(t.Last_verse)

		if err1 != nil || err2 != nil {
			return tsst, errors.New("OH OH ~Specials")
		}

		// TODO: Handle lists of footnotes or crossrefs

	}

	for _, tb := range ms.Tables {
		chapt, err1 := convert_num(tb.Last_chapter)
		vrs, err2 := convert_num(tb.Last_verse)

		if err1 != nil || err2 != nil {
			return tsst, errors.New("OH OH ~Tables")
		}

		table := Table{
			LastChapter: chapt,
			LastVerse:   vrs,
			Table:       tb.Table,
			Additionals: tb.Additionals,
		}

		tsst.Tables = append(tsst.Tables, table)
	}

	return tsst, nil
}

type Migration_Structure struct {
	General      any
	Verses       []Verse_migrate
	Footnotes    []Footnote_migrate
	Crossrefs    []Crossref_migrate
	Tables       []Table_migrate
	Raw_Specials Raw_Special_migrate
	Titles       []Title_migrate
}

type Footnote_migrate struct {
	References         string
	References_chapter string
	References_verse   string
	Text               string
}

type Crossref_migrate struct {
	References         string
	Belongs_to_chapter string
	Belongs_to_verse   string
	Text               string
}

type Raw_Special_migrate struct {
	Kind        string
	Explanation string
	Specials    []Special_migrate
}

type Special_migrate struct {
	Content    string
	Chapter    string
	Last_verse string
	Kind       string
}

type TableRow_migrate map[string][]CellGroup

type CellGroup_migrate map[string]any

type Table_migrate struct {
	Last_chapter string     `json:"last_chapter"`
	Last_verse   string     `json:"last_verse"`
	Table        []TableRow `json:"table"`
	Additionals  []any      `json:"additionals,omitempty"`
}

type Title_migrate struct {
	Kind       string
	Content    string
	Chapter    string
	Last_verse string
	Footnote   []Footnote_migrate
	Crossref   []Crossref_migrate
}

type Verse_migrate struct {
	Global_locator         string
	Chapter                string
	Verse_number           string
	Alternate_verse_number string // carefull
	Text                   string
	Is_a_list_element      bool // carefull
	Position               string
}
