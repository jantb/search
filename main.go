package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os/user"
	"flag"
	"os"
	"path/filepath"
	"time"
	"syscall"
	"github.com/jantb/search/proto"
	"github.com/jantb/search/tail"
	"golang.org/x/net/context"
	"github.com/jantb/search/searchfor"
	"net"
	"google.golang.org/grpc"
	"github.com/jantb/search/searchbox"
)

var filename = flag.String("add", "", "Filename to monitor")
var poll = flag.Bool("poll", false, "use poll")
var db *bolt.DB

const (
	port = ":50051"
)

type server struct{}

func (s *server) Process(ctx context.Context, in *proto.SearchConf) (*proto.SearchRes, error) {
	channel := make(chan proto.SearchRes)
	quitChan := make(chan bool)
	go searchfor.SearchFor(in.Text, int(in.Size_), int64(in.Skipped), channel, quitChan, db)
	r := <-channel
	return &r, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, _ := os.OpenFile("x.txt", os.O_WRONLY | os.O_CREATE | os.O_SYNC, 0755)
	syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)
	flag.Parse()

	go func() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		proto.RegisterSearchServer(s, &server{})
		s.Serve(lis)
	}()

	go func() {
		conn, err := grpc.Dial("localhost" + port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := proto.NewSearchClient(conn)
		time.Sleep(10 * time.Second)
		_, err = c.Process(context.Background(), &proto.SearchConf{Text:[]byte("INFO"), Size_: int64(10), Skipped: int64(0) })

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
	}()

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dbs, err := bolt.Open(filepath.Join(usr.HomeDir, ".search.db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = dbs
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Events"))
		tx.CreateBucketIfNotExists([]byte("Files"))
		tx.CreateBucketIfNotExists([]byte("Meta"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		fi, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}

		if !fi.IsDir() {
			err = db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Files"))
				dir, _ := filepath.Abs(filepath.Dir(*filename))
				filep := filepath.Join(dir, filepath.Base(*filename))
				fileMonitor := proto.FileMonitor{
					Path:filep,
					Offset:0,
					Poll: *poll,
				}
				by, err := fileMonitor.Marshal()
				if err != nil {
					log.Fatal(err)
				}
				b.Put([]byte(filep), by)
				return nil
			})
			return
		}
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))
		c := b.Cursor()
		for k, f := c.First(); k != nil; k, f = c.Next() {
			fileMonitor := proto.FileMonitor{}
			fileMonitor.Unmarshal(f)
			if err != nil {
				log.Fatal(err)
			}
			go tail.TailFile(fileMonitor, db)
		}
		return nil
	})

	searchbox.Run(db)

}

