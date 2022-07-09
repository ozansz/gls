package gui

import (
	"github.com/gdamore/tcell/v2"
)

type Shortcut struct {
	Key     string
	Command string
}

const (
	GridTitleColor       = tcell.ColorRed
	TreeViewTitleColor   = tcell.ColorGreen
	FileInfoTitleColor   = tcell.ColorOrange
	DirectoryColor       = tcell.ColorTurquoise
	BorderColor          = tcell.ColorLightGray
	FileInfoAttrColor    = tcell.ColorPowderBlue
	FileInfoValueColor   = tcell.ColorLightSkyBlue
	SearchFormTitleColor = tcell.ColorLightSkyBlue

	FileInfoTabAttrWidth = 20
)

var (
	keyboardShortcuts = []Shortcut{
		{
			Key:     "q",
			Command: "quit",
		},
		{
			Key:     "ESC",
			Command: "quit",
		},
		{
			Key:     "^C",
			Command: "quit",
		},
		{
			Key:     "c",
			Command: "collapse",
		},
		{
			Key:     "e",
			Command: "expand",
		},
		{
			Key:     "s",
			Command: "search",
		},
		{
			Key:     "x",
			Command: "restore",
		},
		{
			Key:     "o",
			Command: "open",
		},
		{
			Key:     "p",
			Command: "open++",
		},
		{
			Key:     "BS",
			Command: "remove",
		},
		{
			Key:     "DEL",
			Command: "remove",
		},
	}
)
