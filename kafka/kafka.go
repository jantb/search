package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jantb/search/logline"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"os"
	"strconv"
)

func KafkaRead(insertLogLinesChan chan logline.LogLine) {
	env, b := os.LookupEnv("KAFKA")
	var cons string
	if b {
		cons = env
	} else {
		cons = "localhost"
	}
	conn, err := kafka.Dial("tcp", fmt.Sprintf("%s:9092", cons))
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var connLeader *kafka.Conn
	connLeader, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer connLeader.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}

		func(p kafka.Partition) {
			// make a new reader that consumes from topic-A, partition 0, at offset 42
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{"localhost:9092"},
				Topic:     p.Topic,
				Partition: p.ID,
				MinBytes:  0,    // 10KB
				MaxBytes:  10e6, // 10MB
			})
			r.SetOffset(0)

			for {
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					break
				}
				l := logline.LogLine{
					Time: m.Time.UnixNano() / 1000000,
				}
				l.SetSystem(fmt.Sprintf("%s %d %d", m.Topic, m.Partition, m.Offset))
				l.SetLevel(string(m.Key))
				indent, err := json.MarshalIndent(string(m.Value), "", "    ")
				if err != nil {
					l.SetBody(string(m.Value))
				} else {
					l.SetBody(string(indent))
				}

				insertLogLinesChan <- l
				//	fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			}

			if err := r.Close(); err != nil {
				log.Fatal("failed to close reader:", err)
			}
		}(p)
	}
	for k := range m {

		fmt.Println(k)
	}
}
