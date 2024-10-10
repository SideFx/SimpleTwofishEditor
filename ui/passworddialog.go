//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// UI password dialog, using Unison library (c) Richard A. Wilkes
// https://github.com/richardwilkes/unison
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"SimpleTwofishEditor/assets"
	"SimpleTwofishEditor/crypto"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"runtime"
)

const (
	PwdSet = iota + 1
	PwdGet
)

const inpTextSize = 200
const obscureRune = 0x2a

var pwdDialog *unison.Dialog
var inpUpper *unison.Field
var inpLower *unison.Field
var okButton *unison.Button
var cancelButton *unison.Button
var dialogMode int

func ShowPasswordDialog(mode int) int {
	var err error
	dialogMode = mode
	pwdDialog, err = newPasswordDialog()
	if err != nil {
		panic(err)
	}
	return pwdDialog.RunModal()
}

func newPasswordDialog() (*unison.Dialog, error) {
	dialog, err := unison.NewDialog(nil, nil, newPasswordMessagePanel(),
		[]*unison.DialogButtonInfo{unison.NewOKButtonInfo(), unison.NewCancelButtonInfo()},
		unison.NotResizableWindowOption())
	if err == nil {
		wnd := dialog.Window()
		if runtime.GOOS == "windows" {
			if len(titleIcons) > 0 {
				wnd.SetTitleIcons(titleIcons)
			}
		}
		if dialogMode == PwdSet {
			wnd.SetTitle(assets.CapPwdSet)
		} else if dialogMode == PwdGet {
			wnd.SetTitle(assets.CapPwdGet)
		}
		okButton = dialog.Button(unison.ModalResponseOK)
		okButton.ClickCallback = func() {
			crypto.Push([]byte(inpUpper.Text()))
			pwdDialog.StopModal(unison.ModalResponseOK)
		}
		cancelButton = dialog.Button(unison.ModalResponseCancel)
		cancelButton.ClickCallback = func() { pwdDialog.StopModal(unison.ModalResponseCancel) }
		okButton.SetEnabled(false)
		return dialog, nil
	}
	return nil, err
}

func newPasswordMessagePanel() *unison.Panel {
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	lblUpper := unison.NewLabel()
	lblUpper.Font = unison.LabelFont
	lblUpper.SetTitle(assets.CapPwdEnter)
	lblLower := unison.NewLabel()
	lblLower.Font = unison.LabelFont
	lblLower.SetTitle(assets.CapPwdVerify)
	inpUpper = unison.NewField()
	inpUpper.Font = unison.FieldFont
	inpUpper.MinimumTextWidth = inpTextSize
	inpUpper.ObscurementRune = obscureRune
	inpLower = unison.NewField()
	inpLower.Font = unison.FieldFont
	inpLower.MinimumTextWidth = inpTextSize
	inpLower.ObscurementRune = obscureRune
	inpUpper.ModifiedCallback = func(before, after *unison.FieldState) {
		inpUpperModifiedCallback(before, after)
	}
	inpLower.ModifiedCallback = func(before, after *unison.FieldState) {
		inpLowerModifiedCallback(before, after)
	}
	panel.SetLayoutData(&unison.FlexLayoutData{
		MinSize: unison.Size{Width: 300},
		HSpan:   1,
		VSpan:   2,
		VAlign:  align.Middle,
	})
	panel.AddChild(lblUpper)
	panel.AddChild(inpUpper)
	if dialogMode == PwdSet {
		panel.AddChild(lblLower)
		panel.AddChild(inpLower)
	}
	panel.Pack()
	return panel
}

func inpUpperModifiedCallback(_, after *unison.FieldState) {
	okButton.SetEnabled(false)
	if dialogMode == PwdSet {
		if after.Text != "" && after.Text == inpLower.Text() {
			okButton.SetEnabled(true)
		}
	} else {
		if dialogMode == PwdGet {
			if after.Text != "" {
				okButton.SetEnabled(true)
			}
		}
	}
}

func inpLowerModifiedCallback(_, after *unison.FieldState) {
	okButton.SetEnabled(false)
	if after.Text != "" && after.Text == inpUpper.Text() {
		okButton.SetEnabled(true)
	}
}
