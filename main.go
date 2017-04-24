package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/jantb/search/tail"
	"github.com/jantb/search/gui"
)
import (
	_ "net/http/pprof"
	"net/http"
	time "time"
	"encoding/json"
)
var filename = flag.String("add", "", "Filename to monitor")
var poll = flag.Bool("poll", false, "use poll")
var db *bolt.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, _ := os.OpenFile("x.txt", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)
	flag.Parse()

	db = getDb()
	defer db.Close()

	if *filename != "" {
		tail.AddFileToTail(*filename, *poll, db)
		return
	}

	tail.TailAllFiles(db)
	go http.ListenAndServe(":8080", http.DefaultServeMux)
	go func() {
		// Grab the initial stats.
		prev := db.Stats()

		for {
			// Wait for 10s.
			time.Sleep(10 * time.Second)

			// Grab the current stats and diff them.
			stats := db.Stats()
			diff := stats.Sub(&prev)

			// Encode stats to JSON and print to STDERR.
			json.NewEncoder(os.Stderr).Encode(diff)

			// Save stats for the next loop.
			prev = stats
		}
	}()
	gui.Run(db)
}

func getDb() *bolt.DB {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dbs, err := bolt.Open(filepath.Join(usr.HomeDir, ".search.db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = dbs

	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Events"))
		tx.CreateBucketIfNotExists([]byte("Files"))
		tx.CreateBucketIfNotExists([]byte("Meta"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
