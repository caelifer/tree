package walker

import (
	"io/ioutil"
	"log"
	"os"

	// Custom packages
	"github.com/caelifer/tree/node"
)

////////////////////////////////////////////////////////////////////////

// Helper type for the filter function
type Filter func(os.FileInfo) bool

// Distinct type used by the walker package to provide an interface for tree-walking routines.
// It is responsible for keeping track of all processed file system nodes, client-provided filtering
// subs and communication channels.
type TreeWalker struct {
	// Internal implementation is hidden
	counter           struct{ nfiles, ndirs uint }
	output            chan *node.Node
	outputFilterChain []Filter
}

// ------------ Public methods ------------

// Constructor. Returns pointer to the TreeWalker object.
func NewTreeWalker() *TreeWalker {
	return &TreeWalker{}
}

// Adds client-provided filter function, that takes an os.FileInfo object as an argument and returns bool,
// to a chain of filters. Processed file system node (os.FileInfo) is considered valid if and only if all
// registered filters are passed (return true).
func (tw *TreeWalker) AddFilter(f Filter) {
	tw.outputFilterChain = append(tw.outputFilterChain, f)
}

// Get count of all directories and files successfully processed by the TreeWalker.
func (tw *TreeWalker) GetCounts() (uint, uint) {
	return tw.counter.ndirs, tw.counter.nfiles
}

// Main interface to the TreeWalker service. Starts depth-first traversal of the file system from the provided
// root directory (dir). Returns a receive-only *Node communication channel for all found valid nodes.
func (tw *TreeWalker) Traverse(dir string) <-chan *node.Node {
	tw.output = make(chan *node.Node)
	// Process the tree, treating dir as root node
	go tw.walk(dir, "", true)

	// Return channel interface
	return tw.output
}

// ------------ Private methods ------------

// Run nodes through the filter chain. Returns slice of valid nodes.
func (tw *TreeWalker) filter(nodes []os.FileInfo) []os.FileInfo {
	// Preallocate array backing validNodes slice
	validNodes := make([]os.FileInfo, 0, len(nodes))

	// Filter nodes based on the custom filters
	for _, node := range nodes {
		valid := true // Default - true
		for _, fn := range tw.outputFilterChain {
			// Check if filter is tripped
			valid = fn(node)
			if !valid {
				break
			}
		}
		// Only if _all_ filters are passed, we add node to the validNodes slice
		if valid {
			validNodes = append(validNodes, node)
		}
	}

	return validNodes
}

// Send Node pointer to the output channel.
func (tw *TreeWalker) emit(n *node.Node) {
	// Update dir and file counters first. No need to lock up since TreeWalker operations are
	// running a separate gorutine.
	if n.IsDir() {
		tw.counter.ndirs++
	} else {
		tw.counter.nfiles++
	}
	tw.output <- n
}

// Walk the file system node
func (tw *TreeWalker) walk(dir, prefix string, isRoot bool) {
	if isRoot { // Process root node
		// When done processing all top-level entries close channel
		defer close(tw.output)

		if info, err := os.Lstat(dir); err == nil {
			tw.processNode(dir, "", "", node.RootNodeMode, info, false)
		} else {
			log.Printf("WARN: failed to %v\n", err)
		}
	} else { // Process the rest of the nodes - recursive clause

		// Read directory entries identified by dir
		dirents, err := ioutil.ReadDir(dir)

		if err == nil {
			// Success. Filter entries and process valid nodes (depth-first)
			entries := tw.filter(dirents)

			// Keep index of the last node in a list
			lastNodeIndex := len(entries) - 1

			for i, n := range entries {
				// Check if last
				last := i == lastNodeIndex

				// Make sure payload is not nil
				if n == nil {
					log.Fatalf("n == nil at entries[%d]", i)
				}

				// Recursively process individual nodes
				var mode node.NodeMode
				if last {
					mode |= node.LastNodeMode
				}
				tw.processNode(n.Name(), prefix, dir, mode, n, last)
			}
		} else {
			log.Printf("WARN: failed to %s\n", err)
		}
	}
}

// Process node
func (tw *TreeWalker) processNode(name, prefix, parent string, mode node.NodeMode, entry os.FileInfo, last bool) {
	// Stop recursion if entry is nil
	if entry != nil {
		// Emit text representation of the node
		tw.emit(node.NewNode(name, parent, prefix, mode, entry))

		prefixPart := "" // Prefix for the root and first level nodes
		relativePath := name

		// If parent is present, make sure to adjust relativePath and prefixPart
		if parent != "" {
			relativePath = parent + string(os.PathSeparator) + relativePath
			prefixPart = "│   " // Prefix for a regular node
		}

		// Set prefix to space-padded string if this was the last entry in directory
		if last {
			prefixPart = "    "
		}

		// Check if directory and walk it
		if entry.IsDir() {
			tw.walk(relativePath, prefix+prefixPart, false)
		}
	}
}
