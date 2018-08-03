// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package flake8

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestParse(t *testing.T) {

	g := Goblin(t)
	g.Describe("Flake8Parser", func() {

		g.It("Should correctly parse the report", func() {

			messages := Parse("fixture/report.txt")

			expectedMessages := []*MessageLine{
				{
					File:      "untitled.py",
					Line:      42,
					Character: 13,
					ErrorCode: "E741",
					Message:   "ambiguous variable name 'I'",
					Raw:       "./untitled.py:42:13: E741 ambiguous variable name 'I'",
				},
				{
					File:      "untitled.py",
					Line:      44,
					Character: 13,
					ErrorCode: "W741",
					Message:   "ambiguous variable name 'I'",
					Raw:       "./untitled.py:44:13: W741 ambiguous variable name 'I'",
				},
				{
					File:      "untitled2.py",
					Line:      45,
					Character: 9,
					ErrorCode: "E741",
					Message:   "ambiguous variable name 'Dummy'",
					Raw:       "./untitled2.py:45:9: E741 ambiguous variable name 'Dummy'",
				},
				{
					File:      "foo.py",
					Line:      126,
					Character: 1,
					ErrorCode: "E302",
					Message:   "expected 2 blank lines, found 0",
					Raw:       "./foo.py:126:1: E302 expected 2 blank lines, found 0",
				},
				{
					File:      "bar.py",
					Line:      166,
					Character: 32,
					ErrorCode: "E231",
					Message:   "missing whitespace after ','",
					Raw:       "./bar.py:166:32: E231 missing whitespace after ','",
				},
			}

			g.Assert(messages).Equal(expectedMessages)
			// g.Assert(messages[0]).Equal(expectedMessages[0])
			// g.Assert(messages[1]).Equal(expectedMessages[1])
			// g.Assert(messages[2]).Equal(expectedMessages[2])
			// g.Assert(messages[3]).Equal(expectedMessages[3])
			// g.Assert(messages[4]).Equal(expectedMessages[4])
		})

	})

}
