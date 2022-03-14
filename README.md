tree [![Go Report Card](https://goreportcard.com/badge/github.com/caelifer/tree)](https://goreportcard.com/report/github.com/caelifer/tree)
====

`tree` is a re-implementation of tree Linux utility using Go language. This implementation tries to be a Go-idiomatic. Only subset of features are implemented. Noticable absence is no ability to output in HTML format. The implementation might be a bit less efficient then C-based, as it always does lstat(2) on each directory or file node processed.

# Installation
```
go get github.com/caelifer/tree
```

# Usage
```
$ tree -h
Usage of tree:
  -F	show decorations like 'ls -F'
  -a	show hidden files
  -checksum
    	print SHA1 checksum for files
  -d	only show directories
  -f	show relative paths
  -i	do not show indentation lines
  -noreport
    	do not display file and directory counts
  -o string
    	stdout(-)|stderr|/dev/null|file (default "-")
```
# Examples

By default, `tree` will examine current directory. Here is a sample output.
```
$ tree
.
├── README.md
├── formatter
│   └── formatter.go
├── node
│   └── node.go
├── tree.go
└── walker
    └── walker.go

3 directories, 5 files
```

Running `tree` with `-checksum` option will force `-i` and `-f` flags and display SHA1 checksum as a first column. Checksum will be calculated for regular files only. Directory entries will be skipped.

```
$ tree -checksum .
ee90047c6959c8e9dafac944595fd41f845a2438 ./README.md
ce24834954211cbc5f4d59ab8422e290e219d1b8 ./formatter/formatter.go
2455ec80bb114e5709be49c9620408eba1264f7e ./node/node.go
90bc93d03dde9125c3f90335cc2caa398f010cc6 ./tree.go
9b24c61d6f4c7543734ad7abdbd8823b1f4d0c06 ./walker/walker.go

3 directories, 5 files
```
