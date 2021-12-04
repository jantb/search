package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

func readFromPipe(insertChan chan string, insertChanJson chan []byte) {
	var buffer bytes.Buffer
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		line, prefix, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		buffer.Write(line)
		for prefix {
			lineNext, p, err := reader.ReadLine()
			prefix = p
			if err != nil && err == io.EOF {
				break
			}
			buffer.Write(lineNext)
		}
		line = buffer.Bytes()
		buffer.Reset()

		if err != nil && err == io.EOF {
			break
		}
		if strings.HasPrefix(strings.TrimSpace(string(line)), "{") {
			insertChanJson <- line
		} else {
			insertChan <- string(line)
		}
	}
}
