package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Errors = true

	config.Producer.Timeout = 5 * time.Second

	return sarama.NewSyncProducer(brokers, config)
}

// DoRequest sends a message and waits for a response from Kafka
func DoRequest[T any](
	producer sarama.SyncProducer,
	consumer sarama.Consumer,
	payload T,
	topicReq string,
	topicResp string) string {
	// Send message
	if err := SendMessage(producer, payload, topicReq); err != nil {
		return "Failed to send message: " + err.Error()
	}

	// Create partition consumer
	partitionConsumer, err := consumer.ConsumePartition(topicResp, 0, sarama.OffsetNewest)
	if err != nil {
		return "Failed to subscribe to topic: " + err.Error()
	}
	if partitionConsumer == nil {
		return "Kafka partitionConsumer is nil"
	}
	defer partitionConsumer.Close()

	// Context with timeout instead of time.After
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var response string
	select {
	case errMsg := <-partitionConsumer.Errors():
		if errMsg != nil {
			response = "Error while receiving message: " + errMsg.Err.Error()
		}
	case msg := <-partitionConsumer.Messages():
		log.Printf("Topic=%s | Message=%s\n", msg.Topic, string(msg.Value))
		response = string(msg.Value)
	case <-ctx.Done():
		response = "Response timeout expired"
		log.Println(response)
	}

	return response
}

// DoServiceRequest subscribes to a topic, processes incoming messages, and sends a response
func DoServiceRequest(
	producer sarama.SyncProducer,
	c sarama.Consumer,
	stopCh <-chan struct{},
	operation func(string) (interface{}, error),
	topicReq string,
	topicResp string,
) {
	consumer, err := c.ConsumePartition(topicReq, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to subscribe to topic %s: %v", topicReq, err)
		return
	}

	go func() {
		defer log.Printf("Stopped listening on topic %s", topicReq)

		for {
			select {
			case errMsg, ok := <-consumer.Errors():
				if !ok {
					log.Printf("Error channel closed for topic %s", topicReq)
					return
				}
				if errMsg != nil {
					log.Printf("Consumer error: %v", errMsg.Err)
				}

			case msg, ok := <-consumer.Messages():
				if !ok {
					log.Printf("Message channel closed for topic %s", topicReq)
					return
				}
				if msg == nil {
					continue
				}

				payload := string(msg.Value)
				log.Printf("Incoming message: Topic=%s | Value=%s", msg.Topic, payload)

				result, err := operation(payload)
				if err != nil {
					log.Printf("Operation failed: %v", err)
					result = err.Error()
				}

				if sendErr := SendMessage(producer, result, topicResp); sendErr != nil {
					log.Printf(" Failed to send response: %v", sendErr)
				}
			}
		}
	}()
}

// SendMessage serializes and sends a message to Kafka
func SendMessage[T any](producer sarama.SyncProducer, payload T, topic string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent: Topic=%s | Partition=%d | Offset=%d", topic, partition, offset)
	return nil
}
