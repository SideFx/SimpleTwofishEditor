//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Assets - strings
//----------------------------------------------------------------------------------------------------------------------

package assets

const (
	CapNew         = "New"
	CapOpen        = "Open"
	CapSave        = "Save"
	CapPassword    = "Password"
	CapLocked      = "Locked"
	CapUnlocked    = "Unlocked"
	CapLockUnlock  = "Lock/Unlock"
	CapPwdGet      = "Password for decryption"
	CapPwdSet      = "Password for encryption"
	CapPwdEnter    = "Enter password"
	CapPwdVerify   = "Verify password"
	CapSaveChanges = "Save changes"
	CapError       = "Error"
	CapCopy        = "Copy"
	CapCut         = "Cut"
	CapPaste       = "Paste"

	TxtAboutSimpleTwofishEditor = "Simple Twofish Editor v1.0\n(w) 2024 by Jan Buchholz"
	TxtAboutDetails             = "Twofish Go port based on Bruce Schneier's\nreference C implementation:\nhttps://www.schneier.com/academic/twofish/"
	TxtAboutUnison              = "\n\nCredits:\nSimple Twofish Editor has been developed using\nRichard Wilkes' Unison library:\nhttps://github.com/richardwilkes/unison" +
		"\n\nGopher image created at:\nhttps://gopherize.me/"

	UnnamedFile      = "No name." + FileExtension
	UnnamedFileNoExt = "No name"
	AppName          = "Simple Twofish Editor"
	FileExtension    = "twofish"

	ErrFileOpen        = "Error opening file."
	ErrFileRead        = "Error reading file."
	ErrFileWrite       = "Error writing file."
	ErrNoZydeco        = "No Twofish Editor file."
	ErrCorruptedZydeco = "File appears to be corrupted."
	ErrDecryptionError = "Decryption failed."
	ErrEncryptionError = "Encryption failed."
	ErrUnableToDecrypt = "Unable to decrypt file. Please check password entered and try again."
	ErrEmptyFile       = "Empty file detected."

	MsgDocumentModified = "Save changes before closing?"
	MsgWantSave         = "If you don't save, your changes will be lost."
)
