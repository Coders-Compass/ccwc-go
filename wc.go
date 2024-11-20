package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const usage = `
Usage:
	wc [flags] [file1 file2...]
OR
	cat [file ...] | wc [flags]

Flags:
	-l: count lines
	-w: count words
	-c: count bytes
	-m: count characters

NOTE:
1. Flags can be combined, e.g. -lcw
2. If no flags are provided, default is -clw
3. When multiple files are provided, the counts are printed for each file
4. If no files are present, input is read from stdin`

type Counts struct {
	lines int
	words int
	bytes int
	runes int
}

func countAll(r io.Reader) (Counts, error) {
	counts := Counts{}
	reader := bufio.NewReader(r)
	inWord := false
	for {
		r, size, err := reader.ReadRune()
		// See if the error was because we reached the end of the file
		if err == io.EOF {
			if inWord {
				counts.words++
			}
			return counts, nil
		}
		// Something else happened
		if err != nil {
			return Counts{}, fmt.Errorf("error reading file: %v", err)
		}

		counts.runes++
		counts.bytes += size

		if r == '\n' {
			counts.lines++
		}

		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			if inWord {
				inWord = false
				counts.words++
			}
		} else {
			inWord = true
		}
	}
}

type commandLineOptions struct {
	lines bool
	words bool
	bytes bool
	chars bool
	files []string
}

func parseArgs() (*commandLineOptions, error) {
	options := &commandLineOptions{}
	arguments := os.Args[1:]

	for i, argument := range arguments {
		if !strings.HasPrefix(argument, "-") {
			// We found a non-flag argument, rest are files
			options.files = arguments[i:]
			break
		}

		// We have a flag (possibly a group of flags)
		for _, flag := range argument[1:] {
			switch flag {
			case 'l':
				options.lines = true
			case 'w':
				options.words = true
			case 'c':
				options.bytes = true
			case 'm':
				options.chars = true
			default:
				fmt.Printf("invalid flag: %c\n", flag)
				fmt.Println(usage)
				os.Exit(1)
			}
		}
	}

	// If no flags are provided, default to counting lines, words, and bytes (-clw)
	if !options.lines && !options.words && !options.bytes && !options.chars {
		options.lines = true
		options.words = true
		options.bytes = true
	}

	return options, nil
}

func printCounts(counts Counts, fileName string, options *commandLineOptions) {
	if options.lines {
		fmt.Printf("%8d", counts.lines)
	}
	if options.words {
		fmt.Printf("%8d", counts.words)
	}
	if options.chars {
		fmt.Printf("%8d", counts.runes)
	}
	if options.bytes {
		fmt.Printf("%8d", counts.bytes)
	}

	if fileName != "" {
		fmt.Printf(" %s", fileName)
	}
	fmt.Println()
}

func main() {
	options, err := parseArgs()
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	// If there are no files, read from stdin
	if len(options.files) == 0 {
		counts, err := countAll(os.Stdin)
		if err != nil {
			fmt.Println("error: ", err)
			fmt.Println(usage)
			os.Exit(1)
		}
		printCounts(counts, "", options)
		return
	}

	for _, fileName := range options.files {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println("error opening file: ", err)
			os.Exit(1)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("error closing file: %v", err)
				os.Exit(1)
			}
		}(file)

		counts, err := countAll(file)
		if err != nil {
			fmt.Println("error: ", err)
			fmt.Println(usage)

			os.Exit(1)
		}
		printCounts(counts, fileName, options)
	}
}
