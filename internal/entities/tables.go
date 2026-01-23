package entities

// TODO: Make table structure more readable, if time allows

// TableRow represents a row in the table, e.g., {"1": [...]}
type TableRow map[string][]CellGroup

// CellGroup represents a group of cells, e.g., {"cell-1": ..., "cell-2": ...}
type CellGroup map[string]any

// TableData represents the entire table structure
type Table struct {
	Chapter            uint64
	Verse_min_range    uint64
	Verse_max_range    uint64
	Verse_min_notation string
	Verse_max_notation string
	Table              []TableRow
	Additionals        []Additional
}

type Additional struct {
	Footnotes Footnote_additional  // always is "footnotes"
	Crossrefs Crossrefs_additional // always is "crossrefs"
}

// Use it carefully
type Footnote_additional map[string][]Footnote

// Use it carefully
type Crossrefs_additional map[string][]Crossref

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
