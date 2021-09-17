package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jantb/search/logline"
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

	consumer, err := sarama.NewConsumer(cons, sarama.NewConfig())
	if err != nil {

	}
	//defer consumer.Close()

	strings, err := consumer.Topics()
	for _, topic := range strings {
		partitions, err := consumer.Partitions(topic)
		if err != nil {

		}
		var consumers []sarama.PartitionConsumer
		for _, partition := range partitions {

			consumePartition, err := consumer.ConsumePartition(topic, partition, 0)
			//defer consumePartition.Close()
			if err != nil {

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
					indent, err := json.MarshalIndent(string(message.Value), "", "    ")
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
