
package main

import (
	"os"
	"bytes"
	"io/ioutil"
	"unicode/utf8"
	"fmt"
)

func consumeQuoted(line []byte) (field []byte, t []byte) {
	inidx, outidx := 0, 0
	
	// line[0] == '"'
	line = line[1:]
	
	for inidx < len(line) {
		r, size := utf8.DecodeRune(line[inidx:])
		
		switch (r) {
		case '"':
			return line[:outidx], line[inidx+size:]
		case '\\':
			inidx += size
			r, size = utf8.DecodeRune(line[inidx:])
			fallthrough
		default:
			inidx += size
			outidx += utf8.EncodeRune(line[outidx:], r)
		}
	}
	
	panic("quoted end of line")
}

func head(line []byte, separator []byte) (h []byte, t []byte) {
	if line[0] == '"' {
		field, line := consumeQuoted(line)
		if !bytes.HasPrefix(line, separator) && len(line) > 0 {
			fmt.Print("|", string(field), "|\n")
			fmt.Print("|", string(line), "|\n")
			panic("bad file")
		}
		
		if len(line) == 0 {
			return field, nil
		}
		
		return field, line[len(separator):]
	}
	
	idx := bytes.Index(line, separator)
	if idx == -1 {
		return line, nil
	}
	
	return line[:idx], line[idx+len(separator):]
}

func splitLine(line []byte, separator []byte) (result [][]byte) {
	for len(line) > 0 {
		var h []byte
		h, line = head(line, separator)
		
		result = append(result, h)
	}
	
	return
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: $s {separator} {filename} {column filename format}\n", os.Args[0])
		os.Exit(1)
	}
	
	separator := []byte(os.Args[1])
	filename := os.Args[2]
	out_format := os.Args[3]
	
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	
	lines := bytes.Split(file, []byte{'\n'})
	
	// Remove extraneous line (\n is a terminator on UNIX, not separator).
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	
	var splitLines [][][]byte
	var mostFields int
	
	for _, l := range lines {
		fields := splitLine(l, separator)
		splitLines = append(splitLines, fields)
		
		if len(fields) > mostFields {
			mostFields = len(fields)
		}
	}
	
	for i := 0; i < mostFields; i++ {
		outfile, err := os.Create(fmt.Sprintf(out_format, i))
		if err != nil {
			panic(err)
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
