package main

import "time"

func toMillis(time time.Time) int64 {
	return time.UnixNano() / 1000000
}
