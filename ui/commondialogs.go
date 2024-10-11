//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// UI dialogs, using Unison library (c) Richard A. Wilkes
// https://github.com/richardwilkes/unison
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"SimpleTwofishEditor/assets"
	"errors"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/unison"
)

func dialogToSaveChanges() int {
	if isModified {
		msgPanel := unison.NewMessagePanel(assets.MsgDocumentModified, assets.MsgWantSave)
		if dialog, err := unison.NewDialog(unison.DefaultDialogTheme.QuestionIcon, unison.DefaultDialogTheme.QuestionIconInk, msgPanel,
			[]*unison.DialogButtonInfo{unison.NewYesButtonInfo(), unison.NewNoButtonInfo(), unison.NewCancelButtonInfo()},
			unison.NotResizableWindowOption()); err != nil {
			errs.Log(err)
		} else {
			wnd := dialog.Window()
			wnd.SetTitle(assets.CapSaveChanges)
			if len(titleIcons) > 0 {
				wnd.SetTitleIcons(titleIcons)
			}
			return dialog.RunModal()
		}
		return unison.ModalResponseCancel
	}
	return unison.ModalResponseDiscard
}

func dialogToDisplaySystemError(primary string, detail error) {
	var msg string
	var err errs.StackError
	if errors.As(detail, &err) {
		errs.Log(detail)
		msg = err.Message()
	} else {
		msg = detail.Error()
	}
	panel := unison.NewMessagePanel(primary, msg)
	if dialog, err := unison.NewDialog(unison.DefaultDialogTheme.ErrorIcon, unison.DefaultDialogTheme.ErrorIconInk, panel,
		[]*unison.DialogButtonInfo{unison.NewOKButtonInfo()}, unison.NotResizableWindowOption()); err != nil {
		errs.Log(err)
	} else {
		wnd := dialog.Window()
		wnd.SetTitle(assets.CapError)
		if len(titleIcons) > 0 {
			wnd.SetTitleIcons(titleIcons)
		}
		dialog.RunModal()
	}
}

func dialogToDisplayErrorMessage(primary string, detail string) {
	panel := unison.NewMessagePanel(primary, detail)
	if dialog, err := unison.NewDialog(unison.DefaultDialogTheme.ErrorIcon, unison.DefaultDialogTheme.ErrorIconInk, panel,
		[]*unison.DialogButtonInfo{unison.NewOKButtonInfo()}, unison.NotResizableWindowOption()); err != nil {
		errs.Log(err)
	} else {
		wnd := dialog.Window()
		wnd.SetTitle(assets.CapError)
		if len(titleIcons) > 0 {
			wnd.SetTitleIcons(titleIcons)
		}
		dialog.RunModal()
	}
}
