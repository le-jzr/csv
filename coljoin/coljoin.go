// Usage: coljoin {separator} {input filename 1} {input filename 2} ...
//
// Creates a CSV files using input files as columns.
// I.e. the first line of output will contain first lines of all input files, in the order they
// were written on the command line, separated by specified separator (which can be any string).
// 
// The fields of the output CSV are quoted if they contain the separator or quotes.
//
// The output CSV is written to the standard output stream.
//
package main

import (
	"os"
	"fmt"
	"bytes"
	"io/ioutil"
)

func needsEscape(field []byte, separator []byte) bool {
	if bytes.Index(field, separator) != -1 {
		return true
	}
	
	if bytes.Index(field, []byte{'"'}) != -1 {
		return true
	}
	
	return false
}

func escape(field []byte) (result []byte) {
	result = []byte{'"'}
	
	for _, b := range field {
		if b == '"' {
			result = append(result, '"', '"')
		} else {
			result = append(result, b)
		}
	}
	
	result = append(result, '"')
	
	return
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s {separator} {input filename 1} {input filename 2} ... \n", os.Args[0])
		os.Exit(1)
	}
	
	separator := []byte(os.Args[1])
	in_filenames := os.Args[2:]
	
	var split_lines [][][]byte	
	
	// Load all inputs into the array.
	
	for i, in_filename := range in_filenames {
		file, err := ioutil.ReadFile(in_filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open file %s: %s\n", in_filename, err)
			os.Exit(1)
		}
		
		lines := bytes.Split(file, []byte{'\n'})
		// Remove extraneous line (\n is a terminator on UNIX, not separator).
		if len(lines[len(lines)-1]) == 0 {
			lines = lines[:len(lines)-1]
		}
		
		for j, line := range lines {
			if j >= len(split_lines) {
				split_lines = append(split_lines, nil)
				for k := 0; k < i; k++ {
					split_lines[j] = append(split_lines[j], nil)
				}
			}
			
			split_lines[j] = append(split_lines[j], line)
		}
	}
	
	// Print all into a single CSV file.
	
	for _, line := range split_lines {
		for j, field := range line {
			if j != 0 {
				os.Stdout.Write(separator)
			}
			
			if needsEscape(field, separator) {
				os.Stdout.Write(escape(field))
			} else {
				os.Stdout.Write(field)
			}
		}
		
		os.Stdout.Write([]byte{'\n'})
	}
}
