package menu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var entityTypes = map[rune]string{
	// Official entity types declared in RFC 1436
	'0': "TextFile",
	'1': "Menu",
	'2': "CSO",
	// NOTE: This might need to be named something else
	'3': "Error",
	'4': "Macintosh",
	'5': "DOS",
	'6': "UUEncoded",
	'7': "IndexServer",
	'8': "Telnet",
	'9': "Binary",
	// NOTE: Duplicate probably shouldn't be stored in its own DirEntity
	'+': "Duplicate",
	'g': "GIF",
	'I': "Image",
	'T': "TN3270",

	// Unofficial standard entity types
	'd': "Doc",
	'h': "HTML",
	'i': "Info",
	'p': "PNG",
	'r': "RTF",
	's': "Sound",
	'P': "PDF",
	'X': "XML",
}

func ValidEntityType(r rune) bool {
	_, valid := entityTypes[r]
	return valid
}

func GetEntityTypeName(r rune) (name string, ok bool) {
	name, ok = entityTypes[r]
	return
}

type DirEntity struct {
	Type     rune
	UserName string
	Selector string
	Hostname string
	Port     int
}

type Menu struct {
	DirEntities []DirEntity
}

func FromString(menuString string) (menu Menu, err error) {
	lines := strings.Split(menuString, "\r\n")

	endOfMenu := false

	for _, line := range lines {

		if len(line) == 0 {
			continue
		}

		if line == "." {
			endOfMenu = true
			break
		}

		entity := DirEntity{}

		// If the first character is not a valid type, skip this item
		_, valid := entityTypes[rune(line[0])]
		if !valid { continue }

		entity.Type = rune(line[0])

		// Remove first character
		line = line[1:]

		fields := strings.Split(line, "\t")
		if len(fields) < 4 {
			// Ignore malformed item
			continue
		}

		port, err := strconv.Atoi(fields[3])
		if err != nil || port < 0 {
			return menu, errors.New("invalid port")
		}

		entity.UserName = fields[0]
		entity.Selector = fields[1]
		entity.Hostname = fields[2]
		entity.Port = port

		menu.DirEntities = append(menu.DirEntities, entity)
	}

	if !endOfMenu {
		err = errors.New("unterminated menu")
	} else {
		err = nil
	}

	return
}

func (m *Menu) Debugln() {
	for _, e := range m.DirEntities {
		fmt.Printf(
			"type: %q, user_name: %s, selector: %s, hostname: %s, port: %d\n",
			e.Type,
			e.UserName,
			e.Selector,
			e.Hostname,
			e.Port,
		)
	}
}
