package display

import (
	"github.com/agent-e11/pher-client/internal/menu"
	"github.com/gdamore/tcell/v2"
)

var DefaultStyleMap = map[string]tcell.Style {
	"menu": tcell.StyleDefault,
}

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

func DisplayMenu(s tcell.Screen, m menu.Menu, offset int, styleMap map[string]tcell.Style) {
	if styleMap == nil {
		styleMap = DefaultStyleMap
	}
	width, height := s.Size()
	// NOTE: This offset/height logic might be off
	for i := offset; i < height + offset; i++ {
		var entity menu.DirEntity
		var entityStr string
		// If the index is out of bounds, then the string is empty
		if i < 0 || i >= len(m.DirEntities) {
			entityStr = ""
		} else {
			entity = m.DirEntities[i]
			entityStr = entity.UserName
		}
		strLen := len(entityStr)

		for col := 4; col < width; col++ {
			var r rune
			if col - 4 >= strLen {
				r = ' ' // Print spaces until the end of the screen
			} else {
				r = rune(entityStr[col - 4])
			}
			s.SetContent(col, i - offset, r, nil, styleMap["menu"])
		}
	}
}
