tree
====

`tree` is a re-implementation of tree Linux utility using Go language. This implementation tries to be a Go-idiomatic. Only subset of features is implemented. Noticable absence is no ability to output in HTML format. The implementation might be a bit less efficient then C-based, as it always does lstat(2) on each directory or file node processed.

# Installation
```
go get github.com/caelifer/tree
```

# Usage
```
$ tree -h
Usage of tree:
  -F=false: show decorations like 'ls -F'
  -a=false: show hidden files
  -d=false: only show directories
  -f=false: show relative paths
  -i=false: do not show indentation lines
  -noreport=false: do not display file and directory counts
```
# Examples

By default, `tree` will examin current directory. Here is a sample output.
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

Running `tree` with `-checksum` option will force `-i` and `-f` flags and display SHA1 checksum as a first column. Checksum will be calculated for regular files only.


```
$ tree -checksum
                                         .
7961c19cdab070d609cbc4109aad4b67fc612de6 ./README.md
                                         ./formatter
bde4ac3c0434c616759d66a28e5138de2de188b6 ./formatter/formatter.go
                                         ./node
432df6adda59a4c79378d547df546dae349db1f3 ./node/node.go
2a1545e5e37c2c821f193f5d92712c68d767cae2 ./tree.go
                                         ./walker
e5d7bcb68d14544cd135f399b8b2ffa141c8223f ./walker/walker.go

3 directories, 5 files
```
