// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"test/kafka"
// 	"time"

// 	"github.com/IBM/sarama"
// )

// // charset use random string
// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// // stringWithCharset return of random string
// func stringWithCharset(length int, charset string) string {
// 	b := make([]byte, length)
// 	for i := range b {
// 		b[i] = charset[seededRand.Intn(len(charset))]
// 	}
// 	return string(b)
// }

// func main() {
// 	go func() {
// 		for {
// 			fmt.Println("heatbeat")
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()
// 	conf := sarama.NewConfig()
// 	conf.Consumer.MaxWaitTime = 10 * time.Second
// 	conf.Producer.Return.Successes = true
// 	connectionString := []string{
// 		"kafka-kafka-1:9092",
// 	}

// 	// Make Client
// 	client, err := sarama.NewClient(connectionString, conf)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// defer client.Close()

// 	// Make Topic
// 	// testTopic := sarama.CreateTopicsRequest{
// 	// 	TopicDetails: map[string]*sarama.TopicDetail{
// 	// 		"test": &sarama.TopicDetail{
// 	// 			NumPartitions:     5,
// 	// 			ReplicationFactor: 1,
// 	// 		},
// 	// 	},
// 	// }
// 	// deleteTopic := sarama.DeleteTopicsRequest{
// 	// 	Topics: []string{"test"},
// 	// }
// 	// client.LeastLoadedBroker().DeleteTopics(&deleteTopic)
// 	// time.Sleep(1 * time.Second)
// 	// client.LeastLoadedBroker().CreateTopics(&testTopic)
// 	// topics, err := client.Topics()
// 	// fmt.Print(topics)
// 	// fmt.Println(err)

// 	// Make Producer
// 	producer, err := sarama.NewSyncProducerFromClient(client)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// defer producer.Close()

// 	// Make Consumer
// 	consumer, err := sarama.NewConsumerFromClient(client)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// defer consumer.Close()

// 	// Make KafkaChannel
// 	kafkaChannel := kafka.KafkaChannel{
// 		Producer: producer,
// 		Consumer: consumer,
// 		Done:     make(chan bool),
// 	}

// 	// Test Start
// 	go kafkaChannel.Consume("test")
// 	time.Sleep(2 * time.Second)
// 	go func() {
// 		for i := 0; i < 10; i++ {
// 			kafkaChannel.Produce("test", stringWithCharset(10, charset))
// 			// time.Sleep(1 * time.Second)
// 		}
// 	}()

// 	// time.Sleep(10 * time.Second)
// 	<-kafkaChannel.Done
// }

// // func main() {
// // 	testGoRoutin()
// // }

// // func testGoRoutin() {
// // 	var wg sync.WaitGroup
// // 	for i := 0; i < 10; i++ {
// // 		wg.Add(1)
// // 		go func(i int, wg *sync.WaitGroup) {
// // 			time.Sleep(time.Duration(i) * time.Second)
// // 			fmt.Println(i)
// // 			wg.Done()
// // 		}(i, &wg)
// // 	}
// // 	wg.Wait()
// // }

// // func main() {
// // 	app := fiber.New()

// // 	app.Use("/ws", func(c *fiber.Ctx) error {
// // 		// IsWebSocketUpgrade returns true if the client
// // 		// requested upgrade to the WebSocket protocol.
// // 		if websocket.IsWebSocketUpgrade(c) {
// // 			c.Locals("allowed", true)
// // 			return c.Next()
// // 		}
// // 		return fiber.ErrUpgradeRequired
// // 	})

// // 	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
// // 		// c.Locals is added to the *websocket.Conn
// // 		log.Println(c.Locals("allowed"))  // true
// // 		log.Println(c.Params("id"))       // 123
// // 		log.Println(c.Query("v"))         // 1.0
// // 		log.Println(c.Cookies("session")) // ""

// // 		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
// // 		var (
// // 			mt  int
// // 			msg []byte
// // 			err error
// // 		)
// // 		for {
// // 			if mt, msg, err = c.ReadMessage(); err != nil {
// // 				log.Println("read:", err)
// // 				break
// // 			}
// // 			log.Printf("recv: %s", msg)

// // 			if err = c.WriteMessage(mt, msg); err != nil {
// // 				log.Println("write:", err)
// // 				break
// // 			}
// // 			log.Println("write:", msg)
// // 			// time.Sleep(3 * time.Second)
// // 		}

// // 	}))

// // 	log.Fatal(app.Listen(":8080"))
// // 	// Access the websocket server: ws://kafka-kafka-1:3000/ws/123?v=1.0
// // 	// https://www.websocket.org/echo.html
// // }
// ProducerTest
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/IBM/sarama"
)

// Reference : https://pkg.go.dev/github.com/shopify/sarama
func SyncWriter(brokerList []string) sarama.SyncProducer {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		fmt.Println("Failed to start Sarama producer:", err)
	}

	return producer
}

func AsyncWriter(brokerList []string) sarama.AsyncProducer {
	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	//tlsConfig := createTlsConfiguration()
	//if tlsConfig != nil {
	//    config.Net.TLS.Enable = true
	//    config.Net.TLS.Config = tlsConfig
	//}
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		fmt.Println("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			fmt.Println("Failed to write access log entry:", err)
		}
	}()

	return producer
}

func AsyncProducer() {
	var brokerList []string
	brokerList = append(brokerList, "kafka-kafka-1:9092")

	x := AsyncWriter(brokerList)

	x.Input() <- &sarama.ProducerMessage{
		Topic: "call_topic",
		Key:   sarama.StringEncoder("CallId"),
		Value: sarama.StringEncoder("AgentId"),
	}

	// time.Sleep(1000 * time.Millisecond)
}

func SyncPrducer() {
	var brokerList []string
	brokerList = append(brokerList, "kafka-kafka-1:9092")

	x := SyncWriter(brokerList)

	//returned partition, offset, err
	p, o, err := x.SendMessage(&sarama.ProducerMessage{
		Topic:     "call_topic",
		Key:       sarama.StringEncoder("CallId"),
		Value:     sarama.StringEncoder("AgentId"),
		Partition: int32(0),
		//Offset
		//Metadata
		//Timestamp
		//Headers
	})

	if err != nil {
		fmt.Println("Error producer: ", err)
		os.Exit(1)
	}
	fmt.Printf("[W] partition : %d, offset : %d\n", p, o)
}

func main() {
	for i := 0; i < 10; i++ {
		// SyncPrducer()
		AsyncProducer()
	}
	time.Sleep(5 * time.Second)
}
