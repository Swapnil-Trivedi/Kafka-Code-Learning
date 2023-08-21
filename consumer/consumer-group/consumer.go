package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

const (
	topic = "Multibroker-App"
	group = "Group01"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	topic string
	group string
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %d", string(message.Value), message.Timestamp, message.Topic, message.Partition)
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func connectConsumer(broker []string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	//config.Consumer.Group.InstanceId = group

	//creating new consumer instance
	conn, err := sarama.NewConsumerGroup(broker, group, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// client consume to start consuming messages for client
func ClientConsume(wg *sync.WaitGroup, ctx context.Context, client sarama.ConsumerGroup, consumer *Consumer) {
	defer wg.Done()
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := client.Consume(ctx, strings.Split(consumer.topic, ","), consumer); err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
		consumer.ready = make(chan bool)
	}
}

func main() {
	keepRunning := true
	brokerUrl := []string{"localhost:9991"}
	client, err := connectConsumer(brokerUrl)
	if err != nil {
		log.Errorf("Failed to create consumer client")
	}
	//cancel context
	ctx, cancel := context.WithCancel(context.Background())

	//consumer group
	consumer := Consumer{
		ready: make(chan bool),
		topic: topic,
		group: group,
	}

	//adding wait group to the consumer start
	wg := &sync.WaitGroup{}
	wg.Add(1)

	//starting message consumption from client
	go ClientConsume(wg, ctx, client, &consumer)

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	//keeping the main func running
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}

	cancel()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

}
