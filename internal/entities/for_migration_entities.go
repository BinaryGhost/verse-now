package entities

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func convert_num(num string) (uint64, error) {
	if num == "" {
		return 0, nil
	}

	uinti, err := strconv.ParseUint(num, 10, 64)

	// Some things like footnotes can appear the beginning of chapter,
	// thus no 0 evalutation
	if err != nil {
		err_msg := fmt.Sprintf("num is 0 or is '%s'!", err)
		return 0, errors.New(err_msg)
	}

	return uinti, nil
}

func convert_to_right_footnote(f Footnote_migrate) (Footnote, error) {
	chapt, err := convert_num(f.References_chapter)
	i, err2 := convert_to_intermediate_form(f.References_verse)

	if err != nil {
		err_mgs := fmt.Sprintf("%s, %s - OH OH ~Footnotes", err, err2)
		return Footnote{}, errors.New(err_mgs)
	}

	return Footnote{
		References:         f.References,
		Chapter:            chapt,
		Verse_min_range:    i.min_range,
		Verse_max_range:    i.max_range,
		Verse_min_notation: i.min_notation,
		Verse_max_notation: i.max_notation,
		Text:               f.Text,
	}, nil
}

func convert_to_right_crossref(cr Crossref_migrate) (Crossref, error) {
	chapt, err := convert_num(cr.Belongs_to_chapter)
	i, err2 := convert_to_intermediate_form(cr.Belongs_to_verse)

	if err != nil || err2 != nil {
		return Crossref{}, errors.New("OH OH ~Crossrefs")
	}

	return Crossref{
		References:         cr.References,
		Chapter:            chapt,
		Verse_min_range:    i.min_range,
		Verse_max_range:    i.max_range,
		Verse_min_notation: i.min_notation,
		Verse_max_notation: i.max_notation,
		Text:               cr.Text,
	}, nil
}

// Turns a verse-number like '1b' to [string, '1', 'b']
func splitInto_Number_Notation(text string) []string {
	re := regexp.MustCompile(`^(\d+)([a-zA-Z]+)$`)
	matches := re.FindStringSubmatch(text)

	return matches
}

type intermediate_form struct {
	min_range    uint64
	max_range    uint64
	min_notation string
	max_notation string
}

func convert_to_intermediate_form(verse string) (intermediate_form, error) {
	if strings.Contains(verse, "-") {
		_parts := strings.Split(verse, "-")
		var min_notation, max_notation string

		left, err := strconv.ParseUint(_parts[0], 10, 64)
		if err != nil {
			matches := splitInto_Number_Notation(_parts[0])
			left, err = convert_num(matches[1])

			if len(matches) != 3 || err != nil {
				return intermediate_form{}, errors.New("Left 'number' not number")
			}

			min_notation = matches[2]
		}

		right, err := strconv.ParseUint(_parts[1], 10, 64)
		if err != nil {
			matches := splitInto_Number_Notation(_parts[1])
			right, err = convert_num(matches[1])

			if len(matches) != 3 || err != nil {
				return intermediate_form{}, errors.New("Right 'number' not number")
			}

			max_notation = matches[2]
		}

		return intermediate_form{
			min_range:    left,
			max_range:    right,
			min_notation: min_notation,
			max_notation: max_notation,
		}, nil
	}

	matches := splitInto_Number_Notation(verse)

	if len(matches) == 3 {
		num, err := convert_num(matches[1])
		if err != nil {
			return intermediate_form{}, errors.New("Notation number wrong")
		}
		return intermediate_form{
			min_range:    num,
			max_range:    num,
			min_notation: matches[2],
			max_notation: matches[2],
		}, nil
	}

	num, err := convert_num(verse)
	if err != nil {
		return intermediate_form{}, err
	}

	return intermediate_form{
		min_range:    num,
		max_range:    num,
		min_notation: "",
		max_notation: "",
	}, nil
}

func Migration_to_real_Structure(ms *Migration_Structure) (TranslationStructure, error) {
	tsst := TranslationStructure{}
	tsst.General = ms.General

	for _, v := range ms.Verses {
		chapt, err1 := convert_num(v.Chapter)
		i, err2 := convert_to_intermediate_form(v.Verse_number)

		if err1 != nil {
			err_mgs := fmt.Sprintf("%v, OH OH ~Verse", err1)
			return tsst, errors.New(err_mgs)
		}

		if err2 != nil {
			err_mgs := fmt.Sprintf("%v, OH OH not Verse_number_strct again ~Verse", err2)
			return tsst, errors.New(err_mgs)
		}

		verse := Verse{
			Global_locator:         v.Global_locator,
			Chapter:                chapt,
			Verse_min_range:        i.min_range,
			Verse_max_range:        i.max_range,
			Verse_min_notation:     i.min_notation,
			Verse_max_notation:     i.max_notation,
			Alternate_verse_number: v.Alternate_verse_number,
			Text:                   v.Text,
			Is_a_list_element:      v.Is_a_list_element,
			Position:               v.Position,
		}

		tsst.Verses = append(tsst.Verses, verse)
	}

	for _, f := range ms.Footnotes {
		footnote, err := convert_to_right_footnote(f)
		if err != nil {
			return tsst, err
		}

		tsst.Footnotes = append(tsst.Footnotes, footnote)
	}
	for _, cr := range ms.Crossrefs {
		crossref, err := convert_to_right_crossref(cr)
		if err != nil {
			return tsst, err
		}

		tsst.Cross_references = append(tsst.Cross_references, crossref)
	}

	for _, s := range ms.Raw_Specials.Specials {
		chapt, err1 := convert_num(s.Chapter)
		i, err2 := convert_to_intermediate_form(s.Last_verse)

		if err1 != nil || err2 != nil {
			err_msg := fmt.Sprintf("%s, %s - OH OH ~Specials", err1, err2)
			return tsst, errors.New(err_msg)
		}

		special := Special{
			Kind:               s.Kind,
			Content:            s.Content,
			Chapter:            chapt,
			Verse_min_range:    i.min_range,
			Verse_max_range:    i.max_range,
			Verse_min_notation: i.min_notation,
			Verse_max_notation: i.max_notation,
		}

		tsst.Special_Elems = append(tsst.Special_Elems, special)
	}

	for _, t := range ms.Titles {
		chapt, err1 := convert_num(t.Chapter)
		i, err2 := convert_to_intermediate_form(t.Last_verse)

		if err1 != nil || err2 != nil {
			err_msg := fmt.Sprintf("%s, %s - OH OH ~Titles", err1, err2)
			return tsst, errors.New(err_msg)
		}

		tmp_footnote_array := []Footnote{}
		for _, f := range t.Footnote { // Yes, it is a list, but is singular in it's name
			append_f, err := convert_to_right_footnote(f)
			if err != nil {
				return tsst, err
			}

			tmp_footnote_array = append(tmp_footnote_array, append_f)
		}

		tmp_crossref_array := []Crossref{}
		for _, cr := range t.Crossref { // Same here
			append_cr, err := convert_to_right_crossref(cr)
			if err != nil {
				return tsst, err
			}

			tmp_crossref_array = append(tmp_crossref_array, append_cr)
		}

		title := Title{
			Kind:               t.Kind,
			Content:            t.Content,
			Verse_min_range:    i.min_range,
			Verse_max_range:    i.max_range,
			Verse_min_notation: i.min_notation,
			Verse_max_notation: i.max_notation,
			Chapter:            chapt,
			Footnotes:          tmp_footnote_array,
			Crossrefs:          tmp_crossref_array,
		}

		tsst.Titles = append(tsst.Titles, title)
	}

	for _, tb := range ms.Tables {
		chapt, err1 := convert_num(tb.Last_chapter)
		i, err2 := convert_to_intermediate_form(tb.Last_verse)

		if err1 != nil || err2 != nil {
			return tsst, errors.New("OH OH ~Tables")
		}

		var additionals_real Additional
		if len(tb.Additionals) > 0 {
			tmp_footnote_array := []Footnote{}
			for _, f := range tb.Additionals[0].footnotes["footnotes"] {
				append_f, err := convert_to_right_footnote(f)
				if err != nil {
					return tsst, err
				}

				tmp_footnote_array = append(tmp_footnote_array, append_f)
			}

			tmp_crossref_array := []Crossref{}
			for _, cr := range tb.Additionals[0].crossrefs["crossrefs"] {
				append_cr, err := convert_to_right_crossref(cr)
				if err != nil {
					return tsst, err
				}

				tmp_crossref_array = append(tmp_crossref_array, append_cr)
			}

			f := make(Footnote_additional)
			f["footnotes"] = tmp_footnote_array

			c := make(Crossrefs_additional)
			c["crossrefs"] = tmp_crossref_array

			additionals_real = Additional{
				Footnotes: f,
				Crossrefs: c,
			}
		}

		table := Table{
			Chapter:            chapt,
			Verse_min_range:    i.min_range,
			Verse_max_range:    i.max_range,
			Verse_min_notation: i.min_notation,
			Verse_max_notation: i.max_notation,
			Table:              tb.Table,
			Additionals:        []Additional{additionals_real},
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

type CellGroup_migrate map[string]any

type Table_migrate struct {
	Last_chapter string                `json:"last_chapter"`
	Last_verse   string                `json:"last_verse"`
	Table        []TableRow            `json:"table"` // NOTE: IS THE SAME AS THE REAL ENTITY
	Additionals  []additionals_migrate `json:"additionals,omitempty"`
}
type additionals_migrate struct {
	footnotes map[string][]Footnote_migrate // always is "footnotes"
	crossrefs map[string][]Crossref_migrate // always is "crossrefs"
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
