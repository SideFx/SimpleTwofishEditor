//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Load/save preferences as Json file
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"SimpleTwofishEditor/assets"
	"encoding/json"
	"github.com/richardwilkes/unison"
	"io"
	"os"
	"path/filepath"
)

func savePreferences() {
	var fontname = ""
	var fontsize = editorFontSize
	rect := mainWindow.FrameRect()
	if name, ok := fontNameMenu.Selected(); ok {
		fontname = name
	}
	if size, ok := fontSizeMenu.Selected(); ok {
		fontsize = size
	}
	prefs := preferences{
		WindowRect: rect,
		FontName:   fontname,
		FontSize:   fontsize,
		LastFolder: lastOpenFolder,
	}
	j, err := json.Marshal(prefs)
	if err == nil {
		dir, _ := os.UserConfigDir()
		dir = filepath.Join(dir, assets.AppName)
		_, err := os.Stat(dir)
		if err != nil {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				panic(err)
			}
		}
		fname := filepath.Join(dir, preferencesFileName)
		_ = os.WriteFile(fname, j, 0644)
	}
}

func loadPreferences() preferences {
	var prefs preferences
	dir, err := os.UserConfigDir()
	dir = filepath.Join(dir, assets.AppName)
	fname := filepath.Join(dir, preferencesFileName)
	j, err := os.Open(fname)
	if err == nil {
		byteValue, _ := io.ReadAll(j)
		_ = j.Close()
		_ = json.Unmarshal(byteValue, &prefs)
	}
	if prefs.FontSize == "" {
		prefs.FontSize = editorFontSize
	}
	return prefs
}

type preferences struct {
	WindowRect unison.Rect
	FontName   string
	FontSize   string
	LastFolder string
}

const preferencesFileName = "org.janbuchholz.simpletwofisheditor.json"
