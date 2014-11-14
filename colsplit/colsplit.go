// Usage: colsplit {separator} {filename} {column-filename-format}
//
// Splits columns of a CSV file into a separate file for each column.
// Quoted fields, if any, are unquoted in the process.
//
// Separator can be specified as any string.
//
// Filename is the name of the input file.
//
// Column-filename-format is a format string for fmt.Printf(), which specifies
// the schema for name of the output files. This schema should contain substring '%d'
// exactly once, and this is going to be replaced by the number of the column, starting with 0.
// 
package main

import (
	"os"
	"bytes"
	"io/ioutil"
	"unicode/utf8"
	"fmt"
)

func consumeQuoted(line []byte, ln int) (field []byte, t []byte) {
	inidx, outidx := 0, 0
	
	// line[0] == '"'
	line = line[1:]
	
	for inidx < len(line) {
		r, size := utf8.DecodeRune(line[inidx:])
		
		if r == '"' {
			inidx += size
			r, size = utf8.DecodeRune(line[inidx:])
			
			if r != '"' {
				return line[:outidx], line[inidx:]
			}
		}
		
		inidx += size
		outidx += utf8.EncodeRune(line[outidx:], r)
	}
	
	fmt.Fprintln(os.Stderr, "Quoted end of line on line %d.\n", ln)
	
	return line[:outidx], nil
}

func head(line []byte, ln int, separator []byte) (h []byte, t []byte) {
	if line[0] == '"' {
		field, line := consumeQuoted(line, ln)
		
		if bytes.HasPrefix(line, separator) {
			line = line[len(separator):]
		} else if len(line) > 0 {
			fmt.Fprintln(os.Stderr, "Malformed quotation on line %d.\n", ln)
		}

		return field, line
	}
	
	idx := bytes.Index(line, separator)
	if idx == -1 {
		return line, nil
	}
	
	return line[:idx], line[idx+len(separator):]
}

func splitLine(line []byte, ln int, separator []byte) (result [][]byte) {
	for len(line) > 0 {
		var h []byte
		h, line = head(line, ln, separator)
		
		result = append(result, h)
	}
	
	return
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s {separator} {filename} {column filename format}\n", os.Args[0])
		os.Exit(1)
	}
	
	separator := []byte(os.Args[1])
	filename := os.Args[2]
	out_format := os.Args[3]
	
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open file %s: %s\n", filename, err)
		os.Exit(1)
	}
	
	lines := bytes.Split(file, []byte{'\n'})
	
	// Remove extraneous line (\n is a terminator on UNIX, not separator).
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	
	var splitLines [][][]byte
	var mostFields int
	
	for ln, l := range lines {
		fields := splitLine(l, ln, separator)
		splitLines = append(splitLines, fields)
		
		if len(fields) > mostFields {
			mostFields = len(fields)
		}
	}
	
	for i := 0; i < mostFields; i++ {
		outfilename := fmt.Sprintf(out_format, i)
		outfile, err := os.Create(outfilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create file %s: %s\n", outfilename, err)
			os.Exit(1)
		}
		
		for _, l := range splitLines {
			if len(l) > i {
				outfile.Write(l[i])
			}
			outfile.Write([]byte{'\n'})
		}
		
		outfile.Close()
	}
}
