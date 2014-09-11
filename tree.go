package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	// Custom packages
	"github.com/caelifer/tree/formatter"
	"github.com/caelifer/tree/walker"
)

////////////////////////////////////////////////////////////////////////////////
// Globals
////////////////////////////////////////////////////////////////////////////////
// Command line options
var (
	showHidden       = flag.Bool("a", false, "show hidden files")
	showOnlyDirs     = flag.Bool("d", false, "only show directories")
	showDecorations  = flag.Bool("F", false, "show decorations like 'ls -F'")
	showRelativePath = flag.Bool("f", false, "show relative paths")
	hidePrefix       = flag.Bool("i", false, "do not show indentation lines")
	hideCount        = flag.Bool("noreport", false, "do not display file and directory counts")
)

////////////////////////////////////////////////////////////////////////////////
// Start the program
////////////////////////////////////////////////////////////////////////////////
func main() {
	// Parse command line options and parameters
	flag.Parse()

	// ----------------------------------------------------------
	// ------------------   Set node filters   ------------------
	// ----------------------------------------------------------
	// Create a new TreeWalker object
	tw := walker.NewTreeWalker()

	// Add custom output filters
	if !*showHidden {
		// Hide nodes starting with '.'
		tw.AddFilter(func(info os.FileInfo) bool { return []byte(info.Name())[0] != '.' })
	}

	if *showOnlyDirs {
		// Only show directories
		tw.AddFilter(func(info os.FileInfo) bool { return info.IsDir() })
	}

	// ----------------------------------------------------------
	// ------------------ Set formatting rules ------------------
	// ----------------------------------------------------------
	// Create our formatter and set display rules
	format := formatter.NewFormatter()

	// Full path?
	format.SetShowFullPath(*showRelativePath)

	// Hide prefix?
	format.SetShowPrefix(!*hidePrefix)

	// Show decoration?
	format.SetShowDecoration(*showDecorations)

	// Show symlink target?
	format.SetShowSymlinkTarget(true)

	// ----------------------------------------------------------
	// Start processing on the background and get the channel back
	// ----------------------------------------------------------
	// Figure out what directory to traverse
	rootDir := "." // Default - current directory
	if args := flag.Args(); len(args) > 0 {
		rootDir = args[0] // ... or first command line argument
	}

	// Get time
	t0 := time.Now()

	// Display output using io.Reader interface
	io.Copy(os.Stdout, format.NewReader(tw.Traverse(rootDir)))

	// Display node count
	if !*hideCount {
		dcnt, fcnt := tw.GetCounts()
		fmt.Printf("\n%d directories", dcnt)

		// Only display file count if not zero
		if fcnt > 0 {
			fmt.Printf(", %d files", fcnt)
		}

		if os.Getenv("DEBUG") != "" {
			fmt.Printf(" [%s]", time.Since(t0))
		}

		fmt.Printf("\n")
	}
}

// vim: :ts=4:sw=4:noexpandtab:ai:
