package main

import (
	"fmt"
	"os"

	"github.com/IBM/sarama"
)

type KafkaChannel struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
	done     chan bool
}

func main() {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	connectionString := []string{
		"kafka-kafka-1:9092",
	}
	conn, err := sarama.NewClient(connectionString, conf)
	if err != nil {
		fmt.Println("Error client: ", err)
		os.Exit(1)
	}

	producer, err := sarama.NewSyncProducerFromClient(conn)
	if err != nil {
		fmt.Println("Error producer: ", err)
		os.Exit(1)
	}
	consumer, err := sarama.NewConsumerFromClient(conn)
	if err != nil {
		fmt.Println("Error consumer: ", err)
		os.Exit(1)
	}
	kChannel := KafkaChannel{
		Producer: producer,
		Consumer: consumer,
		done:     make(chan bool),
	}
	go func() {
		for i := 0; i < 10; i++ {
			kChannel.Produce("kafka-msg", fmt.Sprintf("Message %d", i))
		}
	}()
	go kChannel.Consume("kafka-msg")
	<-kChannel.done
}

func (k *KafkaChannel) Produce(topic string, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := k.Producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error producer: ", err)
		os.Exit(1)
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
}

func (k *KafkaChannel) Consume(topic string) {
	partitionList, err := k.Consumer.Partitions(topic)
	if err != nil {
		fmt.Println("Error consumer: ", err)
		os.Exit(1)
	}
	for partition := range partitionList {
		pc, err := k.Consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("Error consumer partition: ", err)
			os.Exit(1)
		}
		defer pc.AsyncClose()
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition: %d, Offset: %d, Key: %s, Value: %s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	<-k.done
}
