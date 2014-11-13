
package main

import (
	"os"
	"fmt"
	"bytes"
	"io/ioutil"
)

func escape(field []byte) (result []byte) {
	result = []byte{'"'}
	
	for _, b := range field {
		switch b {
		case '"':
			result = append(result, '\\', '"')
		case '\\':
			result = append(result, '\\', '\\')
		default:
			result = append(result, b)
		}
	}
	
	result = append(result, '"')
	
	return
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: $s {separator} {input filename 1} {input filename 2} ... \n", os.Args[0])
		os.Exit(1)
	}
	
	separator := []byte(os.Args[1])
	in_filenames := os.Args[2:]
	
	var split_lines [][][]byte	
	
	// Load all inputs into the array.
	
	for i, in_filename := range in_filenames {
		file, err := ioutil.ReadFile(in_filename)
		if err != nil {
			panic(err)
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
			
			if bytes.Index(field, separator) == -1 {
				os.Stdout.Write(field)
			} else {
				os.Stdout.Write(escape(field))
			}
		}
		
		os.Stdout.Write([]byte{'\n'})
	}
}
