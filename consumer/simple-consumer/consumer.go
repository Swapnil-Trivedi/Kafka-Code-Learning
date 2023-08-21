package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

const (
	topic = "Multibroker-App"
)

func connectConsumer(broker []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	//config.Consumer.Group.InstanceId = group

	//creating new consumer instance
	conn, err := sarama.NewConsumer(broker, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func StartConsumer(consumer sarama.PartitionConsumer, sigChan chan os.Signal, done chan struct{}, msgCount *int) {
	for {
		select {
		case err := <-consumer.Errors():
			fmt.Println(err)
		case msg := <-consumer.Messages():
			*msgCount++
			fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", *msgCount, string(msg.Topic), string(msg.Value))
		case <-sigChan:
			fmt.Println("Interrupt is detected")
			done <- struct{}{}
		}
	}
}

func main() {
	brokerUrl := []string{"localhost:9991"}
	consumer, err := connectConsumer(brokerUrl)
	if err != nil {
		log.Errorf("Failed to create consumer")
	}
	log.Infof("Consumer up | Topic : %v ", topic)
	//creating topic consumer
	partConsumer, err := consumer.ConsumePartition(topic, 0, 0)
	if err != nil {
		log.Errorf("Unable to consume from topic")
	}
	//channel for interrupt
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	//total message count
	msgCount := 0
	// Get signal for finish
	doneCh := make(chan struct{})
	//starting consumer
	go StartConsumer(partConsumer, sigchan, doneCh, &msgCount)
	<-doneCh
	log.Infof("Total messages processed : %v", msgCount)
	if err := consumer.Close(); err != nil {
		log.Errorf("Consumer closed, cannot send to channel")
	}

}
