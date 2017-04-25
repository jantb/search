package proto

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/boltdb/bolt"
)

// itob returns an 8-byte big endian representation of v.
func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// itob returns an 8-byte big endian representation of v.
func Btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func (e *FileMonitor) Store(db *bolt.DB) {
	by, err := e.Marshal()
	if err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))
		bp := tx.Bucket([]byte("FilesPathToId"))

		if bp.Get([]byte(e.Path)) == nil {
			id, _ := b.NextSequence()
			if err != nil {
				log.Fatal(err)
			}
			b.Put(Itob(id), by)

			bp.Put([]byte(e.Path), Itob(id))
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func GetKeyToPath(path string, db *bolt.DB) (id []byte) {

	err := db.View(func(tx *bolt.Tx) error {
		bp := tx.Bucket([]byte("FilesPathToId"))
		id = bp.Get([]byte(path))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return id
}
func GetPathFromId(id []byte, db *bolt.DB) (path string) {

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))
		var buffer bytes.Buffer
		buffer.Write(b.Get(id))
		var f FileMonitor
		f.Unmarshal(buffer.Bytes())
		path = f.Path
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return path
}
