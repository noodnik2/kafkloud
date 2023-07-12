package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/configs"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/controller"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/handlers"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/kclient"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/server"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/util"
)

type appState struct {
	config *configs.AppConfig
	ctx    context.Context
}

func main() {

	as := &appState{
		ctx:    context.Background(),
		config: configs.Load(),
	}

	countingController := &controller.Counter{}
	errorHandler := newErrorHandler()

	kConsumerFactory := func(processor kclient.Processor, topics []string) *kclient.KafkaConsumer {
		return &kclient.KafkaConsumer{
			Url:                   as.config.KafkaUrl,
			GroupId:               as.config.KafkaGroupId,
			Topics:                topics,
			WaitTimeout:           as.config.KafkaWaitTimeout,
			Ctx:                   as.ctx,
			ComponentErrorHandler: errorHandler,
			Processor:             processor,
			DoneChan:              make(chan time.Time),
		}
	}

	// Launch Kafka consumer
	log.Printf("launching Kafka listener...\n")
	kConsumer := kConsumerFactory(countingController, as.config.KafkaTopics)
	if errKcLaunch := kConsumer.Launch(); errKcLaunch != nil {
		log.Fatal(errKcLaunch)
		//notreached
	}

	// Launch Server
	s := &server.Server{
		Addr:         as.config.ServerConfig.Addr,
		WriteTimeout: as.config.ServerConfig.WriteTimeout,
		ReadTimeout:  as.config.ServerConfig.ReadTimeout,
		IdleTimeout:  as.config.ServerConfig.IdleTimeout,
		Context:      as.ctx,
		Handlers: &handlers.Handlers{
			Counter:          countingController,
			KConsumerFactory: kConsumerFactory,
		},
	}
	log.Printf("launching HTTP server...\n")
	if errServer := s.Launch(errorHandler); errServer != nil {
		log.Fatal(errServer)
		//notreached
	}

	log.Printf("HTTP server running...\n")
	as.gracefulShutdown([]io.Closer{kConsumer, s})
	//notreached
}

func (as *appState) gracefulShutdown(cleanupFns []io.Closer) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan

	log.Printf("terminating due to signal(%v)\n", sig)
	ctx, cancel := context.WithTimeout(as.ctx, as.config.ShutdownWait)

	go func() {
		log.Printf("cleanup started\n")
		defer cancel()
		for _, cleanup := range cleanupFns {
			if cleanup != nil {
				_ = cleanup.Close()
			}
		}
		log.Printf("cleanup finished\n")
	}()

	log.Println("waiting for clean up...")
	<-ctx.Done()

	log.Println("exiting")
	os.Exit(0)
}

func newErrorHandler() util.ComponentErrorHandler {
	var errorCount int
	var lastErrorTime time.Time
	const MaxErrorCount = 1000 // max # of errors allowed within
	const MaxErrorAccumulationWindow = 60 * time.Second
	return func(err error) bool {

		errorCount++
		errorTime := time.Now()
		log.Printf("error %d/%d received(%v)\n", errorCount, MaxErrorCount, err)
		if !lastErrorTime.IsZero() && errorTime.Sub(lastErrorTime) > MaxErrorAccumulationWindow {
			log.Printf("reset error counter from %d to 0\n", errorCount)
			errorCount = 0
		}

		lastErrorTime = errorTime
		if errorCount < MaxErrorCount {
			return false
		}

		log.Printf("terminating process due to unacceptably high frequency of errors\n")
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		return true
	}
}
