package main

import (
	"strconv"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

func ConnectProducer(broker []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	connection, err := sarama.NewSyncProducer(broker, config)
	if err != nil {
		log.Errorf("Error connecting to Kafka broker %v", err.Error())
	}
	return connection, nil
}

func ProduceMessage(topic string, producer sarama.SyncProducer) {
	for i := 11; i <= 21; i++ {
		val := "This is message number : " + strconv.Itoa(i)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(val),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Errorf("Error sending the message %v", err)
		}
		log.Infof("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	}
}

func main() {
	brokerUrl := []string{"localhost:9991"}
	topic := "Multibroker-App"
	producer, err := ConnectProducer(brokerUrl)
	if err != nil {
		log.Errorf("Error failed in generating producer")
	}
	defer producer.Close()
	ProduceMessage(topic, producer)
}
