package gui

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
	"strconv"
	"strings"
)

type Shortcut struct {
	Key     string
	Command string
}

const configurationFile = ".glsrc"

var (
	GridTitleColor       = tcell.ColorRed
	TreeViewTitleColor   = tcell.ColorGreen
	FileInfoTitleColor   = tcell.ColorOrange
	DirectoryColor       = tcell.ColorTurquoise
	BorderColor          = tcell.ColorLightGray
	FileInfoAttrColor    = tcell.ColorPowderBlue
	FileInfoValueColor   = tcell.ColorLightSkyBlue
	SearchFormTitleColor = tcell.ColorLightSkyBlue
	UnmarkedFileColor    = tcell.ColorWhite
	MarkedFileColor      = tcell.ColorRed

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
			Key:     "r",
			Command: "regex search",
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
		{
			Key:     "m",
			Command: "mark",
		},
		{
			Key:     "u",
			Command: "unmark all",
		},
		{
			Key:     "n",
			Command: "create new file",
		},
		{
			Key:     "v",
			Command: "open file in VIM",
		},
		{
			Key:     "d",
			Command: "cp/paste marked files and folders",
		},
	}
)

func init() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("HOME directory couldn't get")
	}
	path := dirname + "/" + configurationFile
	err = readGLSRCFile(path)
	if err != nil {
		return
	}
}

// readGLSRCFile reads .glsrc file and updates the variables.
func readGLSRCFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	customVars := strings.Split(string(data), "\n")
	for _, v := range customVars {
		if len(strings.TrimSpace(v)) == 0 {
			continue
		}
		vv := strings.Split(v, "=")
		key, val := strings.TrimSpace(vv[0]), strings.TrimSpace(vv[1])
		if strings.EqualFold(key, "GridTitleColor") {
			GridTitleColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "TreeViewTitleColor") {
			TreeViewTitleColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "FileInfoTitleColor") {
			FileInfoTitleColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "DirectoryColor") {
			DirectoryColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "BorderColor") {
			BorderColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "FileInfoAttrColor") {
			FileInfoAttrColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "FileInfoValueColor") {
			FileInfoValueColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "SearchFormTitleColor") {
			SearchFormTitleColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "UnmarkedFileColor") {
			UnmarkedFileColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "MarkedFileColor") {
			MarkedFileColor = tcell.GetColor(val)
		}
		if strings.EqualFold(key, "FileInfoTabAttrWidth") {
			fileInfoTabAttrWidth, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			FileInfoTabAttrWidth = fileInfoTabAttrWidth
		}
	}
	return nil
}
