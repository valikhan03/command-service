package models

import(
	"github.com/Shopify/sarama"
)

type EventProducer struct{
	producer sarama.SyncProducer
	eventChan <- chan *Event
}

//must be run in goroutine
func (p *EventProducer) SendEvents () {
	for{
		event := <- p.eventChan
		//get topic from configs
		topic := ""
		eventjson, err := event.Marshal()
		if err != nil{

		}
		msg := sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(eventjson),
		}

		_, _, err = p.producer.SendMessage(&msg)
		if err != nil{

		}
	}
}