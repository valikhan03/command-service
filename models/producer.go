package models

import (
	"github.com/Shopify/sarama"
	"log"
)

type EventProducer struct {
	producer  sarama.SyncProducer
	eventChan <-chan *Event
}

func NewEventProducer(producer sarama.SyncProducer, eventChan <-chan *Event) *EventProducer {
	return &EventProducer{
		producer:  producer,
		eventChan: eventChan,
	}
}

// must be run in goroutine
func (p *EventProducer) SendEvents() {
	configs := GetKafkaConfigs()
	for {
		event := <-p.eventChan
		eventjson, err := event.Marshal()
		if err != nil {
			log.Println(err)
		}
		msg := sarama.ProducerMessage{
			Topic: configs.Topic,
			Value: sarama.ByteEncoder(eventjson),
		}
		_, _, err = p.producer.SendMessage(&msg)
		if err != nil {
			log.Println(err)
		}
	}
}
