//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// About dialog, using Unison library (c) Richard A. Wilkes
// https://github.com/richardwilkes/unison
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"SimpleTwofishEditor/assets"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"strings"
)

func AboutDialog(item unison.MenuItem) {
	dialog, err := unison.NewDialog(nil, nil, newAboutPanel(),
		[]*unison.DialogButtonInfo{unison.NewOKButtonInfo()},
		unison.NotResizableWindowOption())
	if err == nil {
		wnd := dialog.Window()
		wnd.SetTitle(item.Title())
		okButton = dialog.Button(unison.ModalResponseOK)
		okButton.ClickCallback = func() {
			dialog.StopModal(unison.ModalResponseOK)
		}
		dialog.RunModal()
	}
}

func newAboutPanel() *unison.Panel {
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	breakTextIntoLabels(panel, assets.TxtAboutSimpleTwofishEditor, unison.LabelFont.Face().Font(10), false)
	img, err := newImageFromBytes()
	if err == nil {
		picture := unison.NewLabel()
		picture.Drawable = img
		picture.SetLayoutData(&unison.FlexLayoutData{
			HSpan:  1,
			VSpan:  1,
			HAlign: align.Middle,
			VAlign: align.Middle,
			HGrab:  true,
			VGrab:  true,
		})
		panel.AddChild(picture)
	}
	breakTextIntoLabels(panel, assets.TxtAboutDetails, unison.LabelFont.Face().Font(10), true)
	breakTextIntoLabels(panel, assets.TxtAboutUnison, unison.LabelFont.Face().Font(10), false)
	panel.SetLayoutData(&unison.FlexLayoutData{
		MinSize: unison.Size{Width: 500},
		HSpan:   1,
		VSpan:   1,
		VAlign:  align.Middle,
	})
	return panel
}

// taken (and slightly modifield) from Unison dialog.go
func breakTextIntoLabels(panel *unison.Panel, text string, font unison.Font, addSpaceAbove bool) {
	if text != "" {
		returns := 0
		for {
			if i := strings.Index(text, "\n"); i != -1 {
				if i == 0 {
					returns++
					text = text[1:]
				} else {
					part := text[:i]
					l := unison.NewLabel()
					l.Font = font
					l.SetTitle(part)
					l.SetLayoutData(&unison.FlexLayoutData{
						HSpan:  1,
						VSpan:  1,
						HAlign: align.Middle,
						VAlign: align.Middle,
						HGrab:  true,
					})
					if returns > 1 || addSpaceAbove {
						addSpaceAbove = false
						l.SetBorder(unison.NewEmptyBorder(unison.Insets{Top: unison.StdHSpacing}))
					}
					panel.AddChild(l)
					text = text[i+1:]
					returns = 1
				}
			} else {
				if text != "" {
					l := unison.NewLabel()
					l.Font = font
					l.SetTitle(text)
					l.SetLayoutData(&unison.FlexLayoutData{
						HSpan:  1,
						VSpan:  1,
						HAlign: align.Middle,
						VAlign: align.Middle,
						HGrab:  true,
					})
					if returns > 1 || addSpaceAbove {
						l.SetBorder(unison.NewEmptyBorder(unison.Insets{Top: unison.StdHSpacing}))
					}
					panel.AddChild(l)
				}
				break
			}
		}
	}
}

func newImageFromBytes() (*unison.Image, error) {
	newImage, err := unison.NewImageFromBytes(assets.Gopher, 1)
	return newImage, err
}
