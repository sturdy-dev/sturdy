package live

import (
	"fmt"
	"getsturdy.com/api/pkg/comments"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComment(t *testing.T) {
	cases := []struct {
		name            string
		line            int
		context         string
		contextStartsAt int
		fileName        string
		expectedNewLine int
	}{
		{
			name: "not-changed",
			line: 19,
			// comment is on "fmt.Println(i, buzz)"
			context:         "\t\tfmt.Println(i, fizz)\n\t} else if i%5 == 0 {\n\t\tfmt.Println(i, buzz)\n\t} else {\n\t\tfmt.Println(i)\n",
			contextStartsAt: 17,
			fileName:        "fizzbuzz_0.go",
			expectedNewLine: 19,
		},
		{
			name:            "slightly-moved-down",
			line:            14,
			context:         "\t\tfmt.Println(i, fizz)\n\t} else if i%5 == 0 {\n\t\tfmt.Println(i, buzz)\n\t} else {\n\t\tfmt.Println(i)\n",
			contextStartsAt: 12,
			fileName:        "fizzbuzz_1.go",
			expectedNewLine: 23,
		},
		{
			name:            "slightly-moved-up",
			line:            27,
			context:         "\t\tfmt.Println(i, fizz)\n\t} else if i%5 == 0 {\n\t\tfmt.Println(i, buzz)\n\t} else {\n\t\tfmt.Println(i)\n",
			contextStartsAt: 25,
			fileName:        "fizzbuzz_1.go",
			expectedNewLine: 23,
		},
		{
			name:            "real-nothing-changed",
			line:            14,
			context:         "\t{\n\t\tfizz := \"fizz\"\n\t\tbuzz := \"buzz\"\n\n\t\tif i%3 == 0 && i%5 == 0 {\n",
			contextStartsAt: 12,
			fileName:        "fizzbuzz_real.go",
			expectedNewLine: 14,
		},
		{
			name: "line-above-gets-replaced",
			line: 10,
			// comment is on "func fizzbuzz(i int) {"
			context:         "\n// do fizzbuzz (this comment gets replaced)\nfunc fizzbuzz(i int) {\n\tfizz := \"fizz\"\n\tbuzz := \"buzz\"\n",
			contextStartsAt: 8,
			fileName:        "fizzbuzz_1.go",
			expectedNewLine: 12,
		},
		{
			name:            "indentation-changed",
			line:            14,
			context:         "\tbuzz := \"buzz\"\n\n\tif i%3 == 0 && i%5 == 0 {\n\t\tfmt.Println(i, fizz+buzz)\n\t} else if i%3 == 0 {\n",
			contextStartsAt: 12,
			fileName:        "fizzbuzz_2.go",
			expectedNewLine: 15,
		},
		{
			name:            "indentation-changed-moved",
			line:            20,
			context:         "\tbuzz := \"buzz\"\n\n\tif i%3 == 0 && i%5 == 0 {\n\t\tfmt.Println(i, fizz+buzz)\n\t} else if i%3 == 0 {\n",
			contextStartsAt: 18,
			fileName:        "fizzbuzz_2.go",
			expectedNewLine: 15,
		},
		{
			name:            "nothing-matches",
			line:            17,
			context:         "// a\n// b\n// c\n// d\n// e\n",
			contextStartsAt: 15,
			fileName:        "fizzbuzz_2.go",
			expectedNewLine: -1,
		},
		{
			name:            "no-change",
			line:            7,
			context:         "\nfunc main() {\n\tfizzbuzz(50)\n}\n\n", // comment is on "fizzbuzz(50)
			contextStartsAt: 5,
			expectedNewLine: 7,
			fileName:        "fizzbuzz_2.go",
		},
		{
			name:            "no-change-empty-rows",
			line:            11,
			context:         "}\n\nfunc fizzbuzz(i int) {\n\tfizz := \"fizz\"\n\tbuzz := \"buzz\"\n",
			contextStartsAt: 9,
			expectedNewLine: 11,
			fileName:        "fizzbuzz_3.go",
		},
		{
			name:            "no-change-no-empty",
			line:            7,
			context:         "func main() {\n\tfor i := 1; i <= 50; i++ {\n\t\tfizzbuzz(i)\n\t}\n}\n",
			contextStartsAt: 5,
			expectedNewLine: 7,
			fileName:        "fizzbuzz_3.go",
		},
		{
			name:            "1-line-context",
			line:            1,
			context:         "AAA\n",
			contextStartsAt: 1,
			expectedNewLine: 1,
			fileName:        "small_1_line.txt",
		},
		{
			name:            "2-line-context-1",
			line:            1,
			context:         "AAA\nBBB\n",
			contextStartsAt: 1,
			expectedNewLine: 1,
			fileName:        "small_2_line.txt",
		},
		{
			name:            "2-line-context-2",
			line:            2,
			context:         "AAA\nBBB\n",
			contextStartsAt: 1,
			expectedNewLine: 2,
			fileName:        "small_2_line.txt",
		},
		{
			name:            "3-line-context",
			line:            2,
			context:         "AAA\nBBB\nCCC\n",
			contextStartsAt: 1,
			expectedNewLine: 2,
			fileName:        "small_3_line.txt",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fileContents, err := ioutil.ReadFile("testdata/" + tc.fileName)
			assert.NoError(t, err)

			comment := comments.Comment{
				LineStart:           tc.line,
				LineEnd:             tc.line, // Multiline comments are not yet supported anywhere
				LineIsNew:           true,
				Context:             &tc.context,
				ContextStartsAtLine: &tc.contextStartsAt,
			}

			newLocation := fuzzyNewLocation(comment, string(fileContents))

			assert.Equal(t, tc.expectedNewLine, newLocation)
		})
	}
}

func TestContextMatchCount(t *testing.T) {
	alphabetRows := strings.Split(`a
b
c
d
e
f
g
h
i`, "\n")

	cases := []struct {
		rows          []string
		rowStartNum   int
		context       []string
		expectedFuzzy int
	}{
		{
			rows:          alphabetRows,
			rowStartNum:   2,
			context:       []string{"c", "d", "e", "f", "g"},
			expectedFuzzy: 5,
		},
		{
			rows:          alphabetRows,
			rowStartNum:   2,
			context:       []string{"g", "f", "e", "d", "c"}, // it's reversed
			expectedFuzzy: 5,
		},
		{
			rows:          alphabetRows,
			rowStartNum:   2,
			context:       []string{"c", "e", "f", "g", "h"}, // as if one extra line was inserted
			expectedFuzzy: 4,
		},
		{
			rows:          alphabetRows,
			rowStartNum:   2,
			context:       []string{"c", "d", "d2", "e", "f"}, // as if one line was removed
			expectedFuzzy: 4,
		},
		{
			rows:          []string{"", "", "", "", "", "", "", "", "", "", ""},
			rowStartNum:   2,
			context:       []string{"", "", "", "", ""}, // as if one line was removed
			expectedFuzzy: 5,                            // at most one fuzzy match per row
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v/%v", tc.rowStartNum, tc.context), func(t *testing.T) {
			fuzzyMatches := contextMatchCount(tc.rows, tc.rowStartNum, tc.context, noop)
			assert.Equal(t, tc.expectedFuzzy, fuzzyMatches)
		})
	}
}
