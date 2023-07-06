package kclient

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/util"
)

/*
	Relevant links:

	- https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md

*/

type Processor interface {
	ProcessMessage(topic, message string) error
}

type KafkaConsumer struct {
	Url         string
	GroupId     string
	Topics      []string
	WaitTimeout time.Duration
	Ctx         context.Context
	Processor   Processor
	util.ComponentErrorHandler
	DoneChan    chan time.Time
	quitRequest atomic.Bool
}

func (kc *KafkaConsumer) Launch() error {

	log.Printf("kc.kConfig.KafkaUrl(%s)\n", kc.Url)
	consumer, errNewConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kc.Url,
		"group.id":          kc.GroupId,
		//"auto.create.topics.enable": true,
		"allow.auto.create.topics": true,
		//"auto.offset.reset": "earliest",
	})
	if errNewConsumer != nil {
		return errNewConsumer
	}

	topics := kc.Topics
	log.Printf("subscribing Kafka listener(%v)\n", topics)
	if errSubscribe := consumer.SubscribeTopics(topics, nil); errSubscribe != nil {
		log.Printf("can't subscribe to Kafka topics(%v) ...\n", topics)
		closeConsumer(consumer)
		return errSubscribe
	}

	go func() {

		var errKafkaMessagePump error
		defer func() { log.Printf("Kafka listener(%v) stopped; err=%v...\n", topics, errKafkaMessagePump) }()
		defer func() { kc.DoneChan <- time.Now(); log.Printf("Kafka listener exited\n") }()
		defer func() { closeConsumer(consumer) }()

		for {

			if kc.quitRequest.Load() {
				log.Printf("Kafka listener(%v) closing...\n", topics)
				break
			}

			log.Printf("Kafka listener(%v) waiting...\n", topics)

			var kafkaMsg *kafka.Message
			kafkaMsg, errKafkaMessagePump = consumer.ReadMessage(kc.WaitTimeout)
			if errKafkaMessagePump != nil {
				if errKafkaMessagePump.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				log.Printf("Kafka listener(%v) encountered an error: %v msg(%v)\n", topics, errKafkaMessagePump, kafkaMsg)
				if kc.ComponentErrorHandler(errKafkaMessagePump) {
					log.Printf("Kafka listener(%v) stopped due to the error: %v msg(%v)\n", topics, errKafkaMessagePump, kafkaMsg)
					break
				}
				continue
			}
			messageText := string(kafkaMsg.Value)
			var messageTopic string
			if kafkaMsg.TopicPartition.Topic != nil {
				messageTopic = *kafkaMsg.TopicPartition.Topic
			}
			log.Printf("Kafka listener(%v) received topic(%s) message: %s\n", topics, messageTopic, messageText)
			if errProcess := kc.Processor.ProcessMessage(messageTopic, messageText); errProcess != nil {
				log.Printf("Kafka listener(%v) stopped due to an error processing the message: %v msg(%v)\n", topics, errProcess, kafkaMsg)
				break
			}
		}
	}()

	log.Printf("launched Kafka listener(%v)\n", topics)
	return nil
}

func (kc *KafkaConsumer) Close() error {
	log.Printf("signaling Kafka listener(%v) to close\n", kc.Topics)
	kc.quitRequest.Store(true)
	log.Printf("waiting for Kafka listener(%v) to close\n", kc.Topics)
	<-kc.DoneChan
	log.Printf("Kafka listener(%v) has closed\n", kc.Topics)
	return nil
}

func closeConsumer(consumer *kafka.Consumer) {
	log.Printf("closing Kafka listener\n")
	if errCloseListener := consumer.Close(); errCloseListener != nil {
		log.Printf("error closing Kafka listener: %v\n", errCloseListener)
	}
}
