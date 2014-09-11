package node

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////
type NodeMode uint8

const (
	// Flag, Mask
	RootNodeMode NodeMode = 1<<iota
	LastNodeMode
)

// Node type
type Node struct {
	name,    parent,    prefix,    mark string
	mode                                NodeMode
	info                                os.FileInfo
}

func (n *Node) String() string {
	return fmt.Sprintf(
			"\n\tname: %#v,\n\tparent: %#v,\n\tprefix: %#v," +
					"\n\tmark: %#v,\n\tmode: {%#v},\n\tinfo: {\n\t\t%#v\n\t}\n",
		n.name, n.parent, n.prefix, n.mark, n.mode, n.info,
	)
}

// Node's methods

// Getters
func (n *Node) Mark() string {
	mark := ""
	if !n.IsRoot() {
		mark = "├── "
		if n.IsLast() {
			mark = "└── " // last node in directory
		}
	}
	return mark
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Parent() string {
	return n.parent
}

func joinPath(parent, name string) string {
	res := name
	if parent != "" {
		res = parent + string(os.PathSeparator) + res
	}
	return res
}

func (n *Node) FullPath() string {
	return joinPath(n.Parent(), n.Name())
}

func (n *Node) Prefix() string {
	return n.prefix
}

func (n *Node) Decoration() string {
	if n.IsDir() {
		return "/"
	}

	if n.IsSymlink() {
		return "@"
	}

	if n.IsSocket() {
		return "="
	}

	if n.IsPipe() {
		return "|"
	}

	if n.IsExecutable() {
		return "*"
	}

	return ""
}

func (n *Node) SymlinkTarget() string {
	if n.IsSymlink() {
		lpath := n.FullPath()
		if rpath, err := os.Readlink(lpath); err == nil {
			// Check if target is valid
			if _, err := os.Stat(lpath); err != nil {
				rpath += " [bad link]"
			}
			return rpath
		}
	}
	return "[not symlink]"
}

// State checkers

func (n *Node) IsRoot() bool {
	return n.mode & RootNodeMode != 0
}

func (n *Node) IsLast() bool {
	return n.mode & LastNodeMode != 0
}

func (n *Node) IsDir() bool {
	return n.info.IsDir()
}

func (n *Node) IsSymlink() bool {
	return n.info.Mode() & os.ModeSymlink != 0
}

func (n *Node) IsSocket() bool {
	return n.info.Mode() & os.ModeSocket != 0
}
func (n *Node) IsPipe() bool {
	return n.info.Mode() & os.ModeNamedPipe != 0
}

func (n *Node) IsExecutable() bool {
	if n.info.Mode().IsRegular() {
		return (n.info.Mode() & os.ModePerm) & 0111 != 0
	}
	return false
}

// Constructor
func NewNode(name, parent, prefix string, mode NodeMode, info os.FileInfo) *Node {
	return &Node{
		name:   name,
		parent: parent,
		prefix: prefix,
		mode:   mode,
		info:   info,
	}
}
