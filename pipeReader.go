package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
)

func readFromPipe(insertChan chan string, insertChanJson chan map[string]interface{}) {
	var buffer bytes.Buffer
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return
	}
	readFormats()
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
		var j map[string]interface{}
		if err := json.Unmarshal([]byte(line), &j); err != nil {
			insertChan <- string(line)
		} else {
			insertChanJson <- j
		}
	}
}
