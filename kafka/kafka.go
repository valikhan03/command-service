package kafka

import (
	"log"

	"github.com/valikhan03/command-service/models"

	"github.com/Shopify/sarama"
)


func NewProducer() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(models.GetKafkaConfigs().Brokers, config)
	if err != nil{
		log.Fatalf("unable to connect kafka producer: %s\n", err.Error())
	}

	return producer
}

