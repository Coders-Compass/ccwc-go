# ccwc - Word Count Clone

A simple clone of the Unix `wc` utility written in Go. Counts lines, words, characters and bytes in text files or standard input.

## Features
- Counts lines (-l), words (-w), characters (-m), and bytes (-c)
- Handles multiple files
- Supports standard input through pipes
- Unicode-aware character counting
- Combines flags (e.g. -lwm)

## Build
```bash
go build -o ccwc . 
```

## Usage
```bash
ccwc [flags] [files...]

# Default behavior (-clw) on a file
ccwc file.txt

# Count only lines and words in multiple files  
ccwc -lw file1.txt file2.txt

# Count characters from standard input
cat file.txt | ccwc -m

# Count lines from standard input
echo "hello\nworld" | ccwc -l
```

### Flags
- `-l` count lines
- `-w` count words
- `-m` count characters (Unicode support)
- `-c` count bytes

Multiple flags can be combined: `-lwm`  
If no flags are specified, defaults to `-clw`

## Examples with Output
```bash
$ echo "hello world" | ccwc
       1       2      12

$ ccwc -lw test.txt
    7145   58164 test.txt

$ ccwc -c README.md
     1110 README.md
```

## License
MIT License
