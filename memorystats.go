package main

import (
	"runtime"
	"time"
)

func memusage() (uint64, time.Duration, int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return bToMb(m.Alloc), time.Duration(m.PauseTotalNs), runtime.NumGoroutine()
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
