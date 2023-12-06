package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	deliveryChannel := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("Transferiu", "teste", producer, []byte("transfer2"), deliveryChannel)
	go DeliveryReport(deliveryChannel) // async - goroutine
	producer.Flush(6000)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "2031-desenvolvendoprodutoreconsumidor-kafka-1:9092",
		"delivery.timeout.ms": "0",
		"acks":                "all",
		"enable.idempotence":  "true",
	}

	p, err := kafka.NewProducer(configMap)

	if err != nil {
		log.Println(err.Error())
	}

	return p
}

func Publish(msg string, topic string, producer *kafka.Producer, key []byte, deliveryChannel chan kafka.Event) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
		Key:            key,
	}

	err := producer.Produce(message, deliveryChannel)

	if err != nil {
		return err
	}

	return nil
}

func DeliveryReport(deliveryChannel chan kafka.Event) {
	for e := range deliveryChannel {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Println("Erro ao enviar mensagem")
			} else {
				log.Println("Mensagem enviada:", ev.TopicPartition)
			}
		}
	}
}
