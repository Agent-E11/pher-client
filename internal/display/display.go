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
	// Fallback to default style map
	if styleMap == nil {
		styleMap = DefaultStyleMap
	}

	width, height := s.Size()
	row := 0
	col := 0
	dirIdx := offset

	// Loop over the rows in the screen
	for {
		if row >= height {
			return
		}
		if dirIdx > height + offset {
			return
		}

		// Get the current menu item
		var entity menu.DirEntity
		if dirIdx < 0 || dirIdx >= len(m.DirEntities) {
			entity = menu.BlankDirEntity
		} else {
			entity = m.DirEntities[dirIdx]
		}

		// Print the user name, wrapping to the next line if needed
		for _, r := range entity.UserName {
			// If the column is off screen, wrap to next line
			if col >= width {
				col = 0
				row++
				if row >= height {
					return
				}

			}
			s.SetContent(col, row, r, nil, styleMap["menu"])

			col++
		}
		// Print spaces until the end of the current line
		for ; col < width; col++ {
			s.SetContent(col, row, ' ', nil, styleMap["menu"])
		}
		dirIdx++

		// Go to start of next line
		row++
		col = 0
	}
}
