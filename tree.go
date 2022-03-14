package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	showHash         = flag.Bool("checksum", false, "print SHA1 checksum for files")
	output           = flag.String("o", "-", "stdout(-)|stderr|/dev/null|file")
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
		tw.AddFilter(func(info os.FileInfo) bool { return []rune(info.Name())[0] != '.' })
	}

	if *showOnlyDirs {
		// Only show directories
		tw.AddFilter(func(info os.FileInfo) bool { return info.IsDir() })
	}

	// ----------------------------------------------------------
	// ------------------ Set formatting rules ------------------
	// ----------------------------------------------------------

	// Show SHA1 checksum - special case
	if *showHash {

		// Explicitly modify formatting rules
		*showRelativePath = true
		*hidePrefix = true
	}

	// Create our formatter and set display rules
	format := formatter.NewFormatter()

	// Show SHA1 hash
	format.SetShowHash(*showHash)

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
	args := append([]string{}, flag.Args()...)
	if len(args) < 1 {
		args = []string{"."} // Default is a current directory
	}

	// Get time
	t0 := time.Now()

	// Capture output writer in the function literal
	puts := func(w io.Writer, roots []string) {
		for i, root := range roots {
			if i > 0 {
				fmt.Fprintf(w, "\n")
			}
			// Traverse all provided roots
			if _, err := io.Copy(w, format.NewReader(tw.Traverse(root))); err != nil {
				log.Fatalf("Failed to write to output: %v", err)
			}
		}
		// Display node count
		if !*hideCount {
			dcnt, fcnt := tw.GetCounts()
			fmt.Fprintf(w, "\n%d directories", dcnt)

			// Only display file count if not zero
			if fcnt > 0 {
				fmt.Fprintf(w, ", %d files", fcnt)
			}

			if os.Getenv("DEBUG") != "" {
				fmt.Fprintf(w, " [%s]", time.Since(t0))
			}

			fmt.Fprintf(w, "\n")
		}
	}
	// Select output writer
	switch *output {
	case "stdout", "-":
		puts(os.Stdout, args)
	case "stderr":
		puts(os.Stderr, args)
	case "/dev/null":
		// Be nice to Windows sufferers
		puts(ioutil.Discard, args)
	default:
		if out, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY, 0666); err == nil {
			defer out.Close()
			puts(out, args)
		} else {
			log.Fatal(err)
		}
	}
}

// vim: :ts=4:sw=4:noexpandtab:ai:
