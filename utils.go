package main

import (
	"regexp"
	"time"
)

func toMillis(time time.Time) int64 {
	return time.UnixNano() / 1000000
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func Strip(str string) string {
	return re.ReplaceAllString(str, "")
}

func Merge(ch1 <-chan LogLine, ch2 <-chan LogLine) <-chan LogLine {

	out := make(chan LogLine)

	go func() {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		for ok1 || ok2 {
			if !ok2 || (ok1 && v1.Time >= v2.Time) {
				out <- v1
				v1, ok1 = <-ch1
			} else {
				out <- v2
				v2, ok2 = <-ch2
			}
		}
		close(out)
	}()
	return out
}
