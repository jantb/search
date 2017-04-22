package proto

import (
	"log"

	"github.com/boltdb/bolt"
	"bytes"
)

func (e *Meta) Retrieve(db *bolt.DB) {
	var by []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Meta"))
		var buffer bytes.Buffer
		buffer.Write(b.Get([]byte("Meta")))
		by = buffer.Bytes()
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(by) != 0 {
		e.Unmarshal(by)
	}
}

func (e *Meta) Store(db *bolt.DB) {
	by, err := e.Marshal()
	if err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Meta"))
		if err != nil {
			log.Fatal(err)
		}
		b.Put([]byte("Meta"), by)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
