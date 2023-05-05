package kclient

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/configs"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/controller"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/util"
)

/*
	Relevant links:

	- https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md

*/

type KafkaConsumer struct {
	kConfig    *configs.KafkaConfig
	ctx        context.Context
	controller controller.Controller
	quitChan   chan string
}

func NewKafkaConsumer(ctx context.Context, kConfig *configs.KafkaConfig, controller controller.Controller) *KafkaConsumer {
	return &KafkaConsumer{
		kConfig:    kConfig,
		controller: controller,
		ctx:        ctx,
	}
}

func (kc *KafkaConsumer) Launch(errHandler util.ComponentErrorHandler) error {

	fmt.Printf("kc.kConfig.KafkaUrl(%s)\n", kc.kConfig.KafkaUrl)
	kafkaListener, errNewConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kc.kConfig.KafkaUrl,
		"group.id":          kc.kConfig.KafkaGroupId,
		//"auto.create.topics.enable": true,
		"allow.auto.create.topics": true,
		//"auto.offset.reset": "earliest",
	})
	if errNewConsumer != nil {
		return errNewConsumer
	}

	topics := kc.kConfig.KafkaTopics
	fmt.Printf("subscribing to(%v)\n", topics)
	if errSubscribe := kafkaListener.SubscribeTopics(topics, nil); errSubscribe != nil {
		log.Printf("can't subscribe to Kafka topics(%v) ...\n", topics)
		closeConsumer(kafkaListener)
		return errSubscribe
	}

	go func() {
		var errKafkaMessagePump error
		defer func() { closeConsumer(kafkaListener) }()
		for {

			select {
			case <-kc.quitChan:
				log.Printf("Kafka listener closing...\n")
				break
			default:
			}

			log.Printf("Kafka listener waiting...\n")
			var kafkaMsg *kafka.Message
			kafkaMsg, errKafkaMessagePump = kafkaListener.ReadMessage(kc.kConfig.KafkaWaitTimeout)
			if errKafkaMessagePump != nil {
				if errKafkaMessagePump.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				log.Printf("Kafka listener encountered an error: %v msg(%v)\n", errKafkaMessagePump, kafkaMsg)
				if errHandler(errKafkaMessagePump) {
					log.Printf("Kafka listener stopped due to the error: %v msg(%v)\n", errKafkaMessagePump, kafkaMsg)
					break
				}
				continue
			}
			messageText := string(kafkaMsg.Value)
			log.Printf("received Kafka message: %s\n", messageText)
			kc.controller.ProcessMessage(messageText)
		}
		log.Printf("Kafka listener stopped; err=%v...\n", errKafkaMessagePump)
	}()

	log.Printf("launched Kafka listener\n")
	return nil
}

func (kc *KafkaConsumer) Close() error {
	kc.quitChan <- "close"
	return nil
}

func closeConsumer(consumer *kafka.Consumer) {
	log.Printf("closing Kafka listener\n")
	if errCloseListener := consumer.Close(); errCloseListener != nil {
		log.Printf("error closing Kafka listener: %v\n", errCloseListener)
	}
}
