package display

import (
	"github.com/gdamore/tcell/v2"
)

func DrawTextWrap(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1

	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			continue
		}
		if r == '\r' {
			col = x1
			continue
		}
		if r == '\t' {
			// HACK:
			col += 4
			continue
		}
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
