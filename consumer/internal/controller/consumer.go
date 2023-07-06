package controller

import "github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/kclient"

type KConsumerFactory func(processor kclient.Processor, topics []string) *kclient.KafkaConsumer
