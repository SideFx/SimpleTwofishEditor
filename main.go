//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Program startup, using Unison UI library (c) Richard A. Wilkes
// https://github.com/richardwilkes/unison
//----------------------------------------------------------------------------------------------------------------------

package main

import (
	"SimpleTwofishEditor/crypto"
	"SimpleTwofishEditor/ui"
	"github.com/richardwilkes/unison"
	"os"
)

func main() {
	if !crypto.SelfTest() {
		os.Exit(255)
	}
	unison.Start(
		unison.StartupFinishedCallback(func() {
			err := ui.NewMainWindow()
			if err != nil {
				panic(err)
			}
		}),
		unison.QuitAfterLastWindowClosedCallback(func() bool {
			return true
		}),
		unison.AllowQuitCallback(func() bool {
			return ui.AllowQuitCallback()
		}),
	)
}
