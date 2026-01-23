package entities

import (
	"encoding/json"
	"fmt"
	"strings"
)

// I know, I know..., but doing this manually and in golang is a different kind of torture

// AI

// TableRow represents a row in the table, e.g., {"1": [...]}
type TableRow map[string][]CellGroup

// CellGroup represents a group of cells, e.g., {"cell-1": ..., "cell-2": ...}
type CellGroup map[string]any

// TableData represents the entire table structure
type Table struct {
	Last_chapter uint64           `bson:"last_chapter"`
	Last_verse   Vrs_number_strct `bson:"last_verse"`
	Table        []TableRow       `bson:"table"`
	Additionals  []Additional     `bson:"additionals,omitempty"`
}

type Additional struct {
	Footnotes Footnote_additional  // always is "footnotes"
	Crossrefs Crossrefs_additional // always is "crossrefs"
}

// Use it carefully
type Footnote_additional map[string][]Footnote

// Use it carefully
type Crossrefs_additional map[string][]Crossref

// String returns a human-readable string representation of the table
func (td Table) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Chapter: %s, Verse: %s\n", td.Last_chapter, td.Last_verse))

	for _, row := range td.Table {
		for rowKey, cellGroups := range row {
			sb.WriteString(fmt.Sprintf("  Row %s:\n", rowKey))
			for _, cellGroup := range cellGroups {
				for cellKey, cellValue := range cellGroup {
					sb.WriteString(fmt.Sprintf("    %s: ", cellKey))
					switch v := cellValue.(type) {
					case string:
						sb.WriteString(fmt.Sprintf("%s\n", v))
					case []any:
						sb.WriteString(fmt.Sprintf("%v\n", v))
					case map[string]any:
						b, _ := json.MarshalIndent(v, "", "  ")
						sb.WriteString(fmt.Sprintf("%s\n", b))
					default:
						sb.WriteString(fmt.Sprintf("%v\n", v))
					}
				}
			}
		}
	}
	return sb.String()
}

// 			{
//             "last_chapter": "7",
//             "last_verse": "5",
//             "table": [
//                 {
//                     "1": [
//                         {
//                             "cell-1": [
//                                 {
//                                     "global_locator": "NEH 7:8",
//                                     "chapter": "7",
//                                     "verse_number": "8",
//                                     "text": "Abusuafo dodow ni:\n"
//                                 },
//
// 								//Can be a verse or
//
//							   "cell-2": string,
//                             "cell-3": [
//                                 {
//                                     "global_locator": "NEH 7:8",
//                                     "chapter": "7",
//                                     "verse_number": "8",
//                                     "text": "Abusuafo dodow ni:\n"
//                                 },
//								   string,
//                             ]
//                         }
//                     ]
//                 },
//				   "additionals": [
//						{
//							"footnotes": [footnote-array]
//						} [or]
//						{
//							"crossrefs": [crossref-array]
//						}
//				    ]
//			...
