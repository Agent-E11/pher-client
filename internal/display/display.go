package display

import (
	"github.com/agent-e11/pher-client/internal/menu"
	"github.com/gdamore/tcell/v2"
)

var StyleUnsupported = tcell.StyleDefault.Foreground(tcell.ColorGray)
var StyleDocument = tcell.StyleDefault.Foreground(tcell.ColorGreen)

var DefaultStyleMap = map[rune]tcell.Style {
	// Blank
	'\x00': tcell.StyleDefault,
	// Text
	'0': StyleDocument,
	// Menu
	'1': tcell.StyleDefault.
		Foreground(tcell.ColorBlue),
	// CSO
	'2': StyleUnsupported,
	// Error
	'3': StyleUnsupported,
	// Macintosh
	'4': StyleUnsupported,
	// DOS
	'5': StyleUnsupported,
	// UUEncoded
	'6': StyleUnsupported,
	// IndexServer
	'7': tcell.StyleDefault.
		Foreground(tcell.ColorYellow),
	// Telnet
	'8': StyleUnsupported,
	// Binary
	//'9': tcell.StyleDefault.
		//Foreground(Binary).
		//Background(Binary),
	// Duplicate
	'+': tcell.StyleDefault,
	// GIF
	//'g': tcell.StyleDefault.
		//Foreground(GIF).
		//Background(GIF),
	// Image
	//'I': tcell.StyleDefault.
		//Foreground(Image).
		//Background(Image),
	// TN3270
	//'T': tcell.StyleDefault.
		//Foreground(TN3270).
		//Background(TN3270),

	// Doc
	'd': StyleDocument,
	// HTML
	'h': StyleDocument,
	// Info
	'i': tcell.StyleDefault,
	// PNG
	'p': StyleUnsupported,
	// RTF
	'r': StyleDocument,
	// Sound
	//'s': tcell.StyleDefault.
		//Foreground(Sound).
		//Background(Sound),
	// PDF
	'P': StyleDocument,
	// XML
	'X': StyleDocument,
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

func DisplayMenu(s tcell.Screen, m menu.Menu, offset int, styleMap map[rune]tcell.Style) {
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

		style, ok := styleMap[entity.Type]
		if !ok {
			// TODO: Change this back to white
			style = tcell.StyleDefault.Foreground(tcell.ColorRed)
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
			s.SetContent(col, row, r, nil, style)

			col++
		}
		// Print spaces until the end of the current line
		for ; col < width; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
		dirIdx++

		// Go to start of next line
		row++
		col = 0
	}
}
