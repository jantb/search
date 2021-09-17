package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/TylerBrock/colorjson"
	"github.com/jantb/search/logline"
	"log"
	"os"
	"strings"
)

func KafkaRead(insertLogLinesChan chan logline.LogLine, quit chan bool) {
	env, b := os.LookupEnv("KAFKA")
	var cons []string
	if b {
		cons = strings.Split(env, ",")
	} else {
		cons = []string{"localhost:9092"}
	}

	config := sarama.NewConfig()
	config.ClientID = "t"
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumer(cons, config)
	if err != nil {
		log.Println(err)
	}
	//defer consumer.Close()

	topics, err := consumer.Topics()
	for _, topic := range topics {
		partitions, err := consumer.Partitions(topic)
		if err != nil {
			log.Println(topic)
			log.Println(err)
			continue
		}
		var consumers []sarama.PartitionConsumer
		for _, partition := range partitions {
			consumePartition, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			//defer consumePartition.Close()
			if err != nil {
				log.Println(topic)
				log.Println(partition)
				log.Println(err)
				continue
			}
			consumers = append(consumers, consumePartition)
			messages := consumePartition.Messages()
			go func(messages <-chan *sarama.ConsumerMessage, topic string, partition int32) {
				for message := range messages {
					l := logline.LogLine{
						Time: message.Timestamp.UnixNano() / 1000000,
					}
					l.SetSystem(fmt.Sprintf("%s %d %d", topic, partition, message.Offset))
					l.SetLevel(string(message.Key))

					var obj map[string]interface{}
					err := json.Unmarshal(message.Value, &obj)
					if err != nil {
						l.SetBody(string(message.Value))
						insertLogLinesChan <- l
						continue
					}

					indent, err := colorjson.Marshal(obj)
					if err != nil {
						l.SetBody(string(message.Value))
					} else {
						l.SetBody(string(indent))
					}

					insertLogLinesChan <- l
				}
			}(messages, topic, partition)
		}
		go func(quit chan bool, consumer sarama.Consumer, consumers []sarama.PartitionConsumer) {
			select {
			case <-quit:
				for _, partitionConsumer := range consumers {
					_ = partitionConsumer.Close()
				}

				_ = consumer.Close()
				return
			}
		}(quit, consumer, consumers)
	}

}
