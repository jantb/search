package proto

import (
	"log"

	"github.com/boltdb/bolt"
)

func (e *FileMonitor) Store(db *bolt.DB) {
	by, err := e.Marshal()
	if err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))

		if b.Get([]byte(e.Path)) == nil {
			if err != nil {
				log.Fatal(err)
			}
			b.Put([]byte(e.Path), by)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
