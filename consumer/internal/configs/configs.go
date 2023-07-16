package configs

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type KafkaConfig struct {
	KafkaUrl         string        `env:"KAFKA_BROKER"`
	KafkaGroupId     string        `env:"KAFKA_GROUP_ID"`
	KafkaTopics      []string      `env:"KAFKA_TOPICS"`
	KafkaWaitTimeout time.Duration `env:"KAFKA_WAIT_TIMEOUT"`
}

type ServerConfig struct {
	Addr         string        `env:"SERVER_ADDR"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT"`
}

type AppConfig struct {
	*KafkaConfig
	*ServerConfig
	ShutdownWait time.Duration `env:"APP_SHUTDOWN_WAIT"`
}

func Load() *AppConfig {

	serverConfig := ServerConfig{
		Addr:         "localhost:8072",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	kafkaConfig := KafkaConfig{
		KafkaUrl:         "uninitialized:0000",
		KafkaGroupId:     "python_example_group_1",
		KafkaTopics:      []string{"uninitialized"},
		KafkaWaitTimeout: 5 * time.Second,
	}

	appConfig := AppConfig{
		ServerConfig: &serverConfig,
		KafkaConfig:  &kafkaConfig,
		ShutdownWait: 20 * time.Second,
	}

	// apply configuration overrides from environment

	printEnv("before dotenv")
	if errGdeLoad := godotenv.Load(); errGdeLoad != nil {
		log.Printf("WARNING: unable to load 'dotenv' environment variables: %v", errGdeLoad)
	}
	printEnv("after dotenv")
	if errEpsc := env.Parse(&serverConfig); errEpsc != nil {
		log.Fatalf("unable to parse Server environment variables: %v", errEpsc)
	}
	if errEpkc := env.Parse(&kafkaConfig); errEpkc != nil {
		log.Fatalf("unable to parse Kafka environment variables: %v", errEpkc)
	}
	if errEpac := env.Parse(&appConfig); errEpac != nil {
		log.Fatalf("unable to parse App environment variables: %v", errEpac)
	}

	return &appConfig
}

func printEnv(when string) {
	var vs = []string{
		"KAFKA_BROKER",
		"KAFKA_GROUP_ID",
		"KAFKA_TOPICS",
		"KAFKA_WAIT_TIMEOUT",
		"SERVER_ADDR",
		"SERVER_WRITE_TIMEOUT",
		"SERVER_READ_TIMEOUT",
		"SERVER_IDLE_TIMEOUT",
		"APP_SHUTDOWN_WAIT",
	}

	log.Printf("environment(%s)\n", when)
	for _, v := range vs {
		vv := os.Getenv(v)
		log.Printf("  - %s: '%s'\n", v, vv)
	}
}
