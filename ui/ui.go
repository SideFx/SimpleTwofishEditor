//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// UI, using Unison library (c) Richard A. Wilkes
// https://github.com/richardwilkes/unison
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"SimpleTwofishEditor/assets"
	"SimpleTwofishEditor/crypto"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/behavior"
	"github.com/richardwilkes/unison/enums/slant"
	"github.com/richardwilkes/unison/enums/spacing"
	"github.com/richardwilkes/unison/enums/weight"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
)

const toolbarFontSize float32 = 8
const editorFontSize = "10"
const toolbuttonX = 18
const toolbuttonY = 18

const (
	FileNewActionID = unison.UserBaseID + iota
	FileOpenActionID
	FileSaveActionID
	EditPasswordActionID
	EditLockActionID
)

const (
	wndMinWidth  float32 = 768
	wndMinHeight float32 = 480
)

var mainWindow *unison.Window
var (
	newBtn      *unison.Button
	openBtn     *unison.Button
	saveBtn     *unison.Button
	passwordBtn *unison.Button
	lockBtn     *unison.Button
	copyBtn     *unison.Button
	cutBtn      *unison.Button
	pasteBtn    *unison.Button
)
var textEditor *unison.Field
var (
	fontNameMenu *unison.PopupMenu[string]
	fontSizeMenu *unison.PopupMenu[string]
)
var fontName = ""
var fontSize = editorFontSize
var isLocked = false
var isModified = false
var lastOpenFolder string
var lastOpenFile = ""

var (
	FileNewAction      *unison.Action
	FileOpenAction     *unison.Action
	FileSaveAction     *unison.Action
	EditPasswordAction *unison.Action
	EditLockAction     *unison.Action
)

func NewMainWindow() error {
	var err error
	mainWindow, err = unison.NewWindow("")
	if err != nil {
		return err
	}
	registerFonts()
	prefs := loadPreferences()
	initMenuHandler()
	installDefaultMenus(mainWindow)
	content := mainWindow.Content()
	content.SetBorder(unison.NewEmptyBorder(unison.NewUniformInsets(5)))
	content.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: 1,
		VSpacing: 5,
	})
	// Create toolbar buttons
	content.AddChild(createToolbarPanel())
	// Create editor
	content.AddChild(createEditorPanel())
	prepareTitleIcon()
	if runtime.GOOS == "windows" {
		if len(titleIcons) > 0 {
			mainWindow.SetTitleIcons(titleIcons)
		}
	}
	mainWindow.Pack()
	// Set MainWindow size & position
	rect := prefs.WindowRect
	if rect.Width < wndMinWidth {
		rect.Width = wndMinWidth
	}
	if rect.Height < wndMinHeight {
		rect.Height = wndMinHeight
	}
	dispRect := unison.PrimaryDisplay().Usable
	if rect.X == 0 || rect.X > dispRect.Width-rect.Width {
		if dispRect.Width > rect.Width {
			rect.X = (dispRect.Width - rect.Width) / 2
		}
	}
	if rect.Y == 0 || rect.Y > dispRect.Height-rect.Height {
		if dispRect.Height > rect.Height {
			rect.Y = (dispRect.Height - rect.Height) / 2
		}
	}
	mainWindow.SetFrameRect(rect)
	// Set last used folder
	lastOpenFolder = prefs.LastFolder
	if lastOpenFolder == "" {
		lastOpenFolder, _ = os.UserHomeDir()
	}
	// Set font family & size
	fontName = prefs.FontName
	fontSize = prefs.FontSize
	if prefs.FontName != "" {
		fontNameMenu.Select(fontName)
	} else {
		if len(registeredFonts) > 0 {
			fontName = registeredFonts[0].Family
			fontNameMenu.Select(fontName)
		}
	}
	if prefs.FontSize != "" {
		fontSizeMenu.Select(fontSize)
	}
	// Install callback routines
	setCallbacks()
	// Set font for editor field
	setEditorFont()
	mainWindow.ToFront()
	textEditor.RequestFocus()
	// Set empty document
	actionNew()
	return nil
}

func createToolbarPanel() *unison.Panel {
	var err error
	var fontSizes = []string{"8", "10", "11", "12", "13", "14", "15", "16", "18", "20"}
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlowLayout{
		HSpacing: 1,
		VSpacing: unison.StdVSpacing,
	})
	//Tool buttons
	newBtn, err = createButton(assets.CapNew, assets.IconNew)
	if err == nil {
		newBtn.SetEnabled(true)
		newBtn.SetFocusable(false)
		panel.AddChild(newBtn)
	}
	openBtn, err = createButton(assets.CapOpen, assets.IconOpen)
	if err == nil {
		openBtn.SetEnabled(true)
		openBtn.SetFocusable(false)
		panel.AddChild(openBtn)
	}
	saveBtn, err = createButton(assets.CapSave, assets.IconSave)
	if err == nil {
		saveBtn.SetEnabled(true)
		saveBtn.SetFocusable(false)
		panel.AddChild(saveBtn)
	}
	createSpacer(3, panel)
	passwordBtn, err = createButton(assets.CapPassword, assets.IconPassword)
	if err == nil {
		passwordBtn.SetEnabled(true)
		passwordBtn.SetFocusable(false)
		panel.AddChild(passwordBtn)
	}
	lockBtn, err = createButton(assets.CapUnlocked, assets.IconUnlocked)
	if err == nil {
		lockBtn.SetEnabled(true)
		lockBtn.SetFocusable(false)
		panel.AddChild(lockBtn)
	}
	createSpacer(10, panel)
	copyBtn, err = createButton(assets.CapCopy, assets.IconCopy)
	if err == nil {
		copyBtn.SetEnabled(true)
		copyBtn.SetFocusable(false)
		panel.AddChild(copyBtn)
	}
	cutBtn, err = createButton(assets.CapCut, assets.IconCut)
	if err == nil {
		cutBtn.SetEnabled(!isLocked)
		cutBtn.SetFocusable(false)
		panel.AddChild(cutBtn)
	}
	pasteBtn, err = createButton(assets.CapPaste, assets.IconPaste)
	if err == nil {
		pasteBtn.SetEnabled(!isLocked)
		pasteBtn.SetFocusable(false)
		panel.AddChild(pasteBtn)
	}
	createSpacer(10, panel)
	//Font names
	fontNameMenu = unison.NewPopupMenu[string]()
	fontNameMenu.Font = unison.LabelFont.Face().Font(toolbarFontSize)
	loadRegisteredFonts()
	panel.AddChild(fontNameMenu)
	createSpacer(5, panel)
	//Font sizes
	fontSizeMenu = unison.NewPopupMenu[string]()
	fontSizeMenu.Font = unison.LabelFont.Face().Font(toolbarFontSize)
	for _, size := range fontSizes {
		fontSizeMenu.AddItem(size)
	}
	panel.AddChild(fontSizeMenu)
	return panel
}

func createEditorPanel() *unison.Panel {
	textEditor = unison.NewMultiLineField()
	textEditor.SetWrap(true)
	textEditor.AutoScroll = false
	_, prefSize, _ := textEditor.Sizes(unison.Size{})
	textEditor.SetFrameRect(unison.Rect{Size: prefSize})
	scroller := unison.NewScrollPanel()
	//Follow: disable hor. scrolling, Fill: enable vert. scrolling
	scroller.SetContent(textEditor, behavior.Follow, behavior.Fill)
	scroller.SetLayoutData(&unison.FlexLayoutData{
		SizeHint: prefSize,
		HAlign:   align.Fill,
		VAlign:   align.Fill,
		HGrab:    true,
		VGrab:    true,
	})
	unison.InstallDefaultFieldBorder(textEditor, scroller)
	scroller.MouseWheelCallback = func(where, delta unison.Point, mod unison.Modifiers) bool {
		b := scroller.DefaultMouseWheel(where, delta, mod)
		if b {
			scroller.Sync()
		}
		return b
	}
	scroller.ScrollRectIntoViewCallback = func(rect unison.Rect) bool {
		b := scroller.DefaultScrollRectIntoView(rect)
		if b {
			scroller.Sync()
		}
		return b
	}
	return scroller.AsPanel()
}

func setCallbacks() {
	newBtn.ClickCallback = func() { fileNew() }
	openBtn.ClickCallback = func() { fileOpen() }
	saveBtn.ClickCallback = func() { fileSave() }
	passwordBtn.ClickCallback = func() { editPassword() }
	lockBtn.ClickCallback = func() { editLock() }
	copyBtn.ClickCallback = func() { editCopy() }
	cutBtn.ClickCallback = func() { editCut() }
	pasteBtn.ClickCallback = func() { editPaste() }
	fontNameMenu.SelectionChangedCallback = func(popup *unison.PopupMenu[string]) {
		fontNameSelected(fontNameMenu)
	}
	fontSizeMenu.SelectionChangedCallback = func(popup *unison.PopupMenu[string]) {
		fontSizeSelected(fontSizeMenu)
	}
	textEditor.KeyDownCallback = func(keyCode unison.KeyCode, mod unison.Modifiers, repeat bool) bool {
		return textEditorKeyDownCallback(keyCode, mod, repeat)
	}
	textEditor.RuneTypedCallback = func(ch rune) bool {
		return textEditorRuneTypedCallback(ch)
	}
	textEditor.RemoveCmdHandler(unison.CutItemID)
	textEditor.InstallCmdHandlers(unison.CutItemID, func(_ any) bool { return textEditorCanCutOverride() }, func(_ any) { textEditor.Cut() })
	textEditor.RemoveCmdHandler(unison.PasteItemID)
	textEditor.InstallCmdHandlers(unison.PasteItemID, func(_ any) bool { return textEditorCanPasteOverride() }, func(_ any) { textEditor.Paste() })
	textEditor.ModifiedCallback = func(before, after *unison.FieldState) {
		textEditorModifiedCallback(before, after)
	}
	mainWindow.MinMaxContentSizeCallback = func() (minSize, maxSize unison.Size) {
		return windowMinMaxResizeCallback()
	}
	mainWindow.AllowCloseCallback = func() bool {
		return mainWindowAllowClose()
	}
	mainWindow.WillCloseCallback = func() {
		mainWindowWillClose()
	}
}

func mainWindowAllowClose() bool {
	if !isModified {
		return true
	}
	answer := dialogToSaveChanges()
	if answer == unison.ModalResponseDiscard {
		isModified = false
		return true
	}
	if answer == unison.ModalResponseOK {
		return actionSave()
	}
	return false
}

func mainWindowWillClose() {
	savePreferences()
}

func fileNew() {
	answer := dialogToSaveChanges()
	if answer == unison.ModalResponseCancel {
		return
	}
	if answer == unison.ModalResponseOK {
		if !actionSave() {
			return
		}
	}
	actionNew()
}

func fileOpen() {
	answer := dialogToSaveChanges()
	if answer == unison.ModalResponseCancel {
		return
	}
	if answer == unison.ModalResponseOK {
		if !actionSave() {
			return
		}
	}
	actionOpen()
}

func fileSave() {
	actionSave()
}

func editPassword() {
	ShowPasswordDialog(PwdSet)
}

func editLock() {
	isLocked = !isLocked
	setLock(isLocked)
}

func setLock(locked bool) {
	var svgcontent string
	isLocked = locked
	if locked {
		lockBtn.SetTitle(assets.CapLocked)
		svgcontent = assets.IconLocked
	} else {
		lockBtn.SetTitle(assets.CapUnlocked)
		svgcontent = assets.IconUnlocked
	}
	cutBtn.SetEnabled(!isLocked)
	pasteBtn.SetEnabled(!isLocked)
	svg, _ := unison.NewSVGFromContentString(svgcontent)
	if svg != nil {
		lockBtn.Drawable = &unison.DrawableSVG{
			SVG:  svg,
			Size: unison.NewSize(toolbuttonX, toolbuttonY),
		}
	}
}

// Popup menu event handlers
func fontNameSelected(popup *unison.PopupMenu[string]) {
	if name, ok := popup.Selected(); ok {
		fontName = name
		setEditorFont()
	}
}

func fontSizeSelected(popup *unison.PopupMenu[string]) {
	if size, ok := popup.Selected(); ok {
		fontSize = size
		setEditorFont()
	}
}

func setEditorFont() {
	if fontName != "" {
		value, err := strconv.ParseFloat(fontSize, 32)
		if err == nil {
			tmp := isModified
			ifont := unison.MatchFontFace(fontName, weight.Regular, spacing.Standard, slant.Upright)
			state := textEditor.GetFieldState()
			textEditor.SetText(textEditor.Text() + " ")
			textEditor.Font = ifont.Font(float32(value))
			textEditor.ApplyFieldState(state)
			textEditor.MarkForRedraw()
			isModified = tmp
		}
	}
}

func editCopy() {
	textEditor.Copy()
}

func editCut() {
	textEditor.Cut()
}

func editPaste() {
	if textEditor.CanPaste() {
		textEditor.Paste()
	}
}

// Editor field event handlers
func textEditorKeyDownCallback(keyCode unison.KeyCode, mod unison.Modifiers, repeat bool) bool {
	if !isLocked {
		return textEditor.DefaultKeyDown(keyCode, mod, repeat)
	}
	return true
}

func textEditorRuneTypedCallback(ch rune) bool {
	if !isLocked {
		return textEditor.DefaultRuneTyped(ch)
	}
	return true
}

func textEditorCanCutOverride() bool {
	if !isLocked {
		return textEditor.CanCut()
	}
	return false
}

func textEditorCanPasteOverride() bool {
	if !isLocked {
		return textEditor.CanPaste()
	}
	return false
}

func textEditorModifiedCallback(before, after *unison.FieldState) {
	isModified = before.Text != after.Text
}

func windowMinMaxResizeCallback() (minSize, maxSize unison.Size) {
	var _min, _max unison.Size
	_min = unison.NewSize(wndMinWidth, wndMinHeight)
	disp := unison.PrimaryDisplay()
	_max.Width = disp.Usable.Width
	_max.Height = disp.Usable.Height
	return _min, _max
}

func AllowQuitCallback() bool {
	mainWindow.AttemptClose()
	return !isModified
}

func actionNew() {
	mainWindow.SetTitle(assets.AppName + " - " + assets.UnnamedFile)
	textEditor.SetText("")
	lastOpenFile = ""
	isModified = false
	setLock(false)
	crypto.Invalidate() //force new password request
}

func actionOpen() {
	var openFile = ""
	var p = ""
	dialog := unison.NewOpenDialog()
	dialog.SetCanChooseDirectories(false)
	dialog.SetAllowsMultipleSelection(false)
	dialog.SetCanChooseFiles(true)
	dialog.SetInitialDirectory(lastOpenFolder)
	dialog.SetAllowedExtensions(assets.FileExtension)
	if dialog.RunModal() == true {
		p = dialog.Path()
		if p != "" {
			lastOpenFolder, openFile = path.Split(p)
			file, err := os.Open(p)
			if err != nil {
				dialogToDisplaySystemError(assets.ErrFileOpen, err)
				return
			}
			payload, err := io.ReadAll(file)
			_ = file.Close()
			if err != nil {
				dialogToDisplaySystemError(assets.ErrFileRead, err)
				return
			}
			if ShowPasswordDialog(PwdGet) == unison.ModalResponseOK {
				clearText, message := crypto.DecryptPayload(payload)
				if message == "" {
					textEditor.SetText(clearText)
					isModified = false
					setLock(true)
					textEditor.SetSelectionToStart()
					lastOpenFile = openFile
					mainWindow.SetTitle(assets.AppName + " - " + lastOpenFile)
				} else {
					dialogToDisplayErrorMessage(assets.ErrDecryptionError, message)
				}
			}
		}
	}
}

func actionSave() bool {
	var saveFile = ""
	var p = ""
	if !crypto.IsValid() {
		if ShowPasswordDialog(PwdSet) != unison.ModalResponseOK {
			return false
		}
	}
	if lastOpenFile == "" || lastOpenFolder == "" {
		dialog := unison.NewSaveDialog()
		dialog.SetInitialFileName(assets.UnnamedFileNoExt)
		dialog.SetInitialDirectory(lastOpenFolder)
		dialog.SetAllowedExtensions(assets.FileExtension)
		if dialog.RunModal() == true {
			p = dialog.Path()
			lastOpenFolder, lastOpenFile = path.Split(p)
		} else {
			return false
		}
	}
	saveFile = path.Join(lastOpenFolder, lastOpenFile)
	cipherText, err := crypto.EncryptPayload([]byte(textEditor.Text()))
	if err != nil {
		dialogToDisplaySystemError(assets.ErrEncryptionError, err)
		return false
	}
	err = os.WriteFile(saveFile, cipherText, 0644)
	if err != nil {
		dialogToDisplaySystemError(assets.ErrFileWrite, err)
		return false
	}
	isModified = false
	mainWindow.SetTitle(assets.AppName + " - " + lastOpenFile)
	return true
}
