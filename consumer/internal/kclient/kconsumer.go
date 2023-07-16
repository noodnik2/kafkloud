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
	log.Printf("(%v) subscribing Kafka listener\n", topics)
	if errSubscribe := consumer.SubscribeTopics(topics, nil); errSubscribe != nil {
		log.Printf("(%v) can't subscribe to Kafka topics: %v\n", topics, errSubscribe)
		closeConsumer(consumer, topics)
		return errSubscribe
	}

	go func() {

		var errKafkaMessagePump error
		defer func() { log.Printf("(%v) Kafka listener stopped; err=%v...\n", topics, errKafkaMessagePump) }()
		defer func() { kc.DoneChan <- time.Now(); log.Printf("(%v) Kafka listener exited\n", topics) }()
		defer func() { closeConsumer(consumer, topics) }()

		for {

			if kc.quitRequest.Load() {
				log.Printf("(%v) Kafka listener closing...\n", topics)
				break
			}

			log.Printf("(%v) Kafka listener waiting...\n", topics)

			var kafkaMsg *kafka.Message
			kafkaMsg, errKafkaMessagePump = consumer.ReadMessage(kc.WaitTimeout)
			if errKafkaMessagePump != nil {
				if errKafkaMessagePump.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				log.Printf("(%v) Kafka listener encountered an error: %v msg(%v)\n", topics, errKafkaMessagePump, kafkaMsg)
				if kc.ComponentErrorHandler(errKafkaMessagePump) {
					log.Printf("(%v) Kafka listener stopped due to the error: %v msg(%v)\n", topics, errKafkaMessagePump, kafkaMsg)
					break
				}
				continue
			}
			messageText := string(kafkaMsg.Value)
			var messageTopic string
			if kafkaMsg.TopicPartition.Topic != nil {
				messageTopic = *kafkaMsg.TopicPartition.Topic
			}
			log.Printf("(%v) Kafka listener received topic(%s) message: %s\n", topics, messageTopic, messageText)
			if errProcess := kc.Processor.ProcessMessage(messageTopic, messageText); errProcess != nil {
				log.Printf("(%v) Kafka listener stopped due to an error processing the message: %v msg(%v)\n", topics, errProcess, kafkaMsg)
				break
			}
		}
	}()

	log.Printf("(%v) launched Kafka listener\n", topics)
	return nil
}

func (kc *KafkaConsumer) Close() error {
	log.Printf("(%v) signaling Kafka listener to close\n", kc.Topics)
	kc.quitRequest.Store(true)
	log.Printf("(%v) waiting for Kafka listener to close\n", kc.Topics)
	<-kc.DoneChan
	log.Printf("(%v) Kafka listener has closed\n", kc.Topics)
	return nil
}

func closeConsumer(consumer *kafka.Consumer, topics []string) {
	log.Printf("(%v) closing Kafka listener\n", topics)
	if errCloseListener := consumer.Close(); errCloseListener != nil {
		log.Printf("(%v) error closing Kafka listener: %v\n", topics, errCloseListener)
	}
}
