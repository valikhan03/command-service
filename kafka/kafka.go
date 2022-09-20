package kafka

import (
	"auctions-service/models"
	"log"

	"github.com/Shopify/sarama"
)


func NewProducer() sarama.SyncProducer {
	config := &sarama.Config{}
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(models.GetKafkaConfigs().Brokers, config)
	if err != nil{
		log.Fatalf("unable to connect kafka producer: %v\n", err)
	}

	return producer
}

