package kafka

import (
	"log"
	"github.com/Shopify/sarama"
)


func NewProducer() sarama.SyncProducer {
	config := &sarama.Config{}
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{}, config)
	if err != nil{
		log.Fatalf("unable to connect kafka producer: %v\n", err)
	}

	return producer
}

