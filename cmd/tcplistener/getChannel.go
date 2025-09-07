package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		defer f.Close()
		currentLine := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("Error: %s\n", err.Error())
				break
			}
			currentLine += string(data[:n])
			lineParts := strings.Split(currentLine, "\n")
			if len(lineParts) > 1 {
				for _, line := range lineParts[:len(lineParts)-1] {
					lines <- line
				}
				currentLine = lineParts[len(lineParts)-1]
			}

		}
	}()
	return lines
}
