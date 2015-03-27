package formatter

import (
	"fmt"
	"io"

	// Local packages
	"github.com/caelifer/tree/node"
)

////////////////////////////////////////////////////////////////////////
type FormatMode uint8

const (
	// Implemented
	ShowFullPathMode FormatMode = 1 << iota
	ShowPrefixMode
	ShowDecorationMode
	ShowSymlinkTargetMode
	ShowHashChecksumMode

	// Future features
	ShowUserGroupMode
	ShowFilePermissionMode
	ShowFileSizeMode
)

////////////////////////////////////////////////////////////////////////
type Formatter struct {
	in   <-chan *node.Node // Receive-only channel of *node.Node(s)
	mode FormatMode
}

func (f *Formatter) String() string {
	return fmt.Sprintf("mode: %08b", f.mode)
}

func (f *Formatter) Next() (string, error) {
	text := ""
	err := io.EOF // EOF is return if channel is closed or in invalid state

	// Get node pointer and channel status
	if n, ok := <-f.in; ok {
		err = nil // Reset error

		// Transform and format n.Node output
		text = n.Name()

		// ------------ Prepends ------------
		// Show SHA1 checksum
		if f.ShowHash() {
			text = n.Checksum()
		}


		// Show relative path
		if f.ShowFullPath() {
			text += " " + n.FullPath()
		}

		// Show prefixes and marking
		if f.ShowPrefix() {
			// Show n's marking
			text = n.Prefix() + n.Mark() + text
		}

		// ------------ Appends ------------

		// Show n decoration, i.e. / for directory, @ for symlink, etc.
		if f.ShowDecoration() {
			text += n.Decoration()
		}

		// Show symlink target
		if f.ShowSymlinkTarget() {
			if n.IsSymlink() {
				text += " → " + n.SymlinkTarget()
			}
		}
	}

	return text, err
}

// Implements io.Reader interface
func (f *Formatter) Read(p []byte) (int, error) {
	var n int // Zero-initialized
	text, err := f.Next()

	if err == nil {
		// Add new line
		n = copy(p, []byte(text+"\n"))
	}

	return n, err
}

// Getters
func (f *Formatter) ShowFullPath() bool {
	return f.mode&ShowFullPathMode != 0
}

func (f *Formatter) ShowPrefix() bool {
	return f.mode&ShowPrefixMode != 0
}

func (f *Formatter) ShowDecoration() bool {
	return f.mode&ShowDecorationMode != 0
}

func (f *Formatter) ShowSymlinkTarget() bool {
	return f.mode&ShowSymlinkTargetMode != 0
}

func (f *Formatter) ShowHash() bool {
	return f.mode&ShowHashChecksumMode != 0
}

// Setters
func (f *Formatter) SetShowFullPath(cond bool) {
	if cond {
		// Set bit
		f.mode |= ShowFullPathMode
	} else {
		// Unset
		f.mode &^= ShowFullPathMode
	}
}

func (f *Formatter) SetShowPrefix(cond bool) {
	if cond {
		// Set bit
		f.mode |= ShowPrefixMode
	} else {
		// Unset
		f.mode &^= ShowPrefixMode
	}
}

func (f *Formatter) SetShowDecoration(cond bool) {
	if cond {
		// Set bit
		f.mode |= ShowDecorationMode
	} else {
		// Unset
		f.mode &^= ShowDecorationMode
	}
}

func (f *Formatter) SetShowSymlinkTarget(cond bool) {
	if cond {
		// Set bit
		f.mode |= ShowSymlinkTargetMode
	} else {
		// Unset
		f.mode &^= ShowSymlinkTargetMode
	}
}

// Display SHA1 checksum for regular file
func (f *Formatter) SetShowHash(cond bool) {
	if cond {
		// Set bit
		f.mode |= ShowHashChecksumMode
	} else {
		// Unset
		f.mode &^= ShowHashChecksumMode
	}
	
}


// Reader wrapper
func (f *Formatter) NewReader(in <-chan *node.Node) io.Reader {
	f.in = in
	return f
}

// Constructor
func NewFormatter() *Formatter {
	return new(Formatter) // Zero-init
}
