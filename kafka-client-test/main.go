package main

import (
	"fmt"

	"github.com/IBM/sarama"
)

func main() {
	var nTntIdx int32 = 0 // Partition Index Set

	config := sarama.NewConfig()
	//config.Version = version

	//default 256
	config.ChannelBufferSize = 1000000
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange

	client, err := sarama.NewClient([]string{"kafka-kafka-1:9092"}, config)

	if err != nil {
		panic(err)
	}

	// if you want completed message -> 0
	lastoffset, err := client.GetOffset("call_topic", nTntIdx, sarama.OffsetNewest)

	if err != nil {
		panic(err)
	}

	// if consumer group isn't exist , create it
	//group, err := sarama.NewConsumerGroupFromClient("consumer-group-name", client)
	consumer, err := sarama.NewConsumerFromClient(client)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//Read Specific Partition From Topic
	partitionConsumer, err := consumer.ConsumePartition("call_topic", nTntIdx, lastoffset)
	//ctx := context.Background()
	//handler := exampleConsumerGroupHandler{}
	//err := group.Consume(ctx, []string{"call_topic"}, handler)

	if err != nil {
		panic(err)
	}

	// if ConsumerGroup Skip Underlines
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	consumed := 0

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Topic %s Consumed message offset %d , Partition %d\n", msg.Topic, msg.Offset, msg.Partition)
			consumed++
			fmt.Printf("Consumed: %d\n", consumed)
			fmt.Println(string(msg.Key))
			fmt.Println(string(msg.Value))
			fmt.Println("")
		}
	}
}
