package kafka

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type KafkaChannel struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
	Done     chan bool
}

func (k *KafkaChannel) Produce(topic string, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(message),
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := k.Producer.SendMessage(msg)
	_ = partition
	_ = offset
	if err != nil {
		fmt.Println("Error producer: ", err)
		os.Exit(1)
	}
	fmt.Printf("[W] Message is stored in topic(%s)/partition(%d)/offset(%d)/msg(%s)\n", topic, partition, offset, message)
}

func messagePrinter(pc sarama.PartitionConsumer, wg *sync.WaitGroup) {
outer:
	for {
		select {
		// case err := <-pc.Errors():
		// 	fmt.Println("Err", err)
		case msg := <-pc.Messages():
			if msg == nil {
				fmt.Println("msg is nil")
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("[R] Partition: %d, Offset: %d, Key: %s, Value: %s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		case <-time.After(5 * time.Second):
			fmt.Println("timeout")
			break outer
		default:
			if pc == nil {
				fmt.Println("pc is nil")
				break outer
			}
		}
	}
	defer func() {
		wg.Done()
		pc.Close()
		fmt.Println("wait group is done")
	}()
}

func (k *KafkaChannel) Consume(topic string) {
	partitionList, err := k.Consumer.Partitions(topic)
	if err != nil {
		fmt.Println("Error consumer: ", err)
		os.Exit(1)
	}
	fmt.Println(partitionList)
	var wg sync.WaitGroup
	for i, partition := range partitionList {
		fmt.Println(i, partition)
		wg.Add(1)
		pc, err := k.Consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("Error consumer partition: ", err)
			os.Exit(1)
		}

		go messagePrinter(pc, &wg)
	}
	wg.Wait()
	fmt.Println("consumer is closed")
	k.Done <- true
}
