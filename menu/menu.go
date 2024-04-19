package menu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type EntityType rune

// TODO: Rename MenuT type, or the struct Menu
// NOTE: Duplicate probably shouldn't be stored in its own DirEntity
// NOTE: This might need to be named something else

// Official entity type declared in RFC 1436
const (
	TextFile    EntityType = '0'
	MenuT       EntityType = '1'
	CSO         EntityType = '2'
	Error       EntityType = '3'
	Macintosh   EntityType = '4'
	DOS         EntityType = '5'
	UUEncoded   EntityType = '6'
	IndexServer EntityType = '7'
	Telnet      EntityType = '8'
	Binary      EntityType = '9'
	Duplicate   EntityType = '+'
	GIF         EntityType = 'g'
	Image       EntityType = 'I'
	TN3270      EntityType = 'T'
)

// Unofficial standard entity types
const (
	Doc   EntityType = 'd'
	HTML  EntityType = 'h'
	Info  EntityType = 'i'
	PNG   EntityType = 'p'
	RTF   EntityType = 'r'
	Sound EntityType = 's'
	PDF   EntityType = 'P'
	XML   EntityType = 'X'
)

type DirEntity struct {
	Type     EntityType
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

		switch line[0] {
		// FIXME: This is a dumb way of doing it, there must be a better way
		// Future me can deal with it though
		case byte(TextFile):
			entity.Type = TextFile
		case byte(MenuT):
			entity.Type = MenuT
		case byte(CSO):
			entity.Type = CSO
		case byte(Error):
			entity.Type = Error
		case byte(Macintosh):
			entity.Type = Macintosh
		case byte(DOS):
			entity.Type = DOS
		case byte(UUEncoded):
			entity.Type = UUEncoded
		case byte(IndexServer):
			entity.Type = IndexServer
		case byte(Telnet):
			entity.Type = Telnet
		case byte(Binary):
			entity.Type = Binary
		case byte(Duplicate):
			entity.Type = Duplicate
		case byte(GIF):
			entity.Type = GIF
		case byte(Image):
			entity.Type = Image
		case byte(TN3270):
			entity.Type = TN3270

		case byte(Doc):
			entity.Type = Doc
		case byte(HTML):
			entity.Type = HTML
		case byte(Info):
			entity.Type = Info
		case byte(PNG):
			entity.Type = PNG
		case byte(RTF):
			entity.Type = RTF
		case byte(Sound):
			entity.Type = Sound
		case byte(PDF):
			entity.Type = PDF
		case byte(XML):
			entity.Type = XML
		default:
			// Ignore malformed type
			continue
		}

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
