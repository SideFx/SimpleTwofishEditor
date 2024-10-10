//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// UI utilities, using Unison library (c) Richard A. Wilkes
// https://github.com/richardwilkes/unison
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"SimpleTwofishEditor/assets"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

var registeredFonts []unison.FontFaceDescriptor
var titleIcons []*unison.Image

func registerFonts() {
	var desc unison.FontFaceDescriptor
	var err error
	desc, err = unison.RegisterFont(assets.FontDejaVu)
	if err == nil {
		registeredFonts = append(registeredFonts, desc)
	}
	desc, err = unison.RegisterFont(assets.FontJetbrains)
	if err == nil {
		registeredFonts = append(registeredFonts, desc)
	}
	desc, err = unison.RegisterFont(assets.FontRobotoMono)
	if err == nil {
		registeredFonts = append(registeredFonts, desc)
	}
	desc, err = unison.RegisterFont(assets.FontSourceCodePro)
	if err == nil {
		registeredFonts = append(registeredFonts, desc)
	}
	desc, err = unison.RegisterFont(assets.FontSpaceMono)
	if err == nil {
		registeredFonts = append(registeredFonts, desc)
	}
}

func loadRegisteredFonts() {
	for _, f := range registeredFonts {
		fontNameMenu.AddItem(f.Family)
	}
}

func createSpacer(width float32, panel *unison.Panel) {
	spacer := &unison.Panel{}
	spacer.Self = spacer
	spacer.SetSizer(func(_ unison.Size) (minSize, prefSize, maxSize unison.Size) {
		minSize.Width = width
		prefSize.Width = width
		maxSize.Width = width
		return
	})
	panel.AddChild(spacer)
}

func newSVGButton(svg *unison.SVG) *unison.Button {
	btn := unison.NewButton()
	btn.HideBase = true
	btn.Drawable = &unison.DrawableSVG{
		SVG:  svg,
		Size: unison.NewSize(toolbuttonX, toolbuttonY),
	}
	btn.Font = unison.LabelFont.Face().Font(toolbarFontSize)
	return btn
}

func createButton(title string, svgcontent string) (*unison.Button, error) {
	svg, err := unison.NewSVGFromContentString(svgcontent)
	if err != nil {
		return nil, err
	}
	btn := newSVGButton(svg)
	btn.SetTitle(title)
	btn.SetLayoutData(align.Middle)
	return btn, nil
}

func installDefaultMenus(wnd *unison.Window) {
	unison.DefaultMenuFactory().BarForWindow(wnd, func(m unison.Menu) {
		unison.InsertStdMenus(m, AboutDialog, nil, nil)
		fileMenu := m.Menu(unison.FileMenuID)
		f := fileMenu.Factory()
		fileMenu.InsertItem(0, FileNewAction.NewMenuItem(f))
		fileMenu.InsertItem(1, FileOpenAction.NewMenuItem(f))
		fileMenu.InsertItem(2, FileSaveAction.NewMenuItem(f))
		fileMenu.InsertSeparator(3, true)
		editMenu := m.Menu(unison.EditMenuID)
		e := editMenu.Factory()
		editMenu.InsertSeparator(-1, true)
		editMenu.InsertItem(-1, EditPasswordAction.NewMenuItem(e))
		editMenu.InsertItem(-1, EditLockAction.NewMenuItem(e))
	})
}

func initMenuHandler() {
	FileNewAction = &unison.Action{
		ID:         FileNewActionID,
		Title:      assets.CapNew,
		KeyBinding: unison.KeyBinding{KeyCode: unison.KeyN, Modifiers: unison.OSMenuCmdModifier()},
		ExecuteCallback: func(_ *unison.Action, _ any) {
			fileNew()
		},
	}
	FileOpenAction = &unison.Action{
		ID:         FileOpenActionID,
		Title:      assets.CapOpen,
		KeyBinding: unison.KeyBinding{KeyCode: unison.KeyO, Modifiers: unison.OSMenuCmdModifier()},
		ExecuteCallback: func(_ *unison.Action, _ any) {
			fileOpen()
		},
	}
	FileSaveAction = &unison.Action{
		ID:         FileSaveActionID,
		Title:      assets.CapSave,
		KeyBinding: unison.KeyBinding{KeyCode: unison.KeyS, Modifiers: unison.OSMenuCmdModifier()},
		ExecuteCallback: func(_ *unison.Action, _ any) {
			fileSave()
		},
	}
	EditPasswordAction = &unison.Action{
		ID:         EditPasswordActionID,
		Title:      assets.CapPassword,
		KeyBinding: unison.KeyBinding{KeyCode: unison.KeyP, Modifiers: unison.OSMenuCmdModifier()},
		ExecuteCallback: func(_ *unison.Action, _ any) {
			editPassword()
		},
	}
	EditLockAction = &unison.Action{
		ID:         EditLockActionID,
		Title:      assets.CapLockUnlock,
		KeyBinding: unison.KeyBinding{KeyCode: unison.KeyL, Modifiers: unison.OSMenuCmdModifier()},
		ExecuteCallback: func(_ *unison.Action, _ any) {
			lockBtn.Click()
		},
	}
}

func prepareTitleIcon() {
	newImage, err := unison.NewImageFromBytes(assets.Lock, 1)
	if err == nil {
		titleIcons = append(titleIcons, newImage)
	}
}

func newImageFromBytes() (*unison.Image, error) {
	newImage, err := unison.NewImageFromBytes(assets.Gopher, 1)
	return newImage, err
}
