package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/configs"
    "github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/controller"
    "github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/kclient"
    "github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/server"
    "github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/util"
)

type appState struct {
    config     *configs.AppConfig
    ctx        context.Context
    controller controller.Controller
}

func main() {

    as := &appState{
        ctx:        context.Background(),
        config:     configs.Load(),
        controller: controller.NewController(),
    }

    componentErrorHandler := newErrorHandler()

    // Launch Kafka consumer
    log.Printf("launching Kafka listener...\n")
    kc := kclient.NewKafkaConsumer(as.ctx, as.config.KafkaConfig, as.controller)
    if errKcLaunch := kc.Launch(componentErrorHandler); errKcLaunch != nil {
        log.Fatal(errKcLaunch)
        //notreached
    }

    // Launch Server
    s := server.NewServer(as.ctx, as.config.ServerConfig, as.controller)
    log.Printf("launching HTTP server...\n")
    if errServer := s.Launch(componentErrorHandler); errServer != nil {
        log.Fatal(errServer)
        //notreached
    }

    log.Printf("HTTP server running...\n")
    as.gracefulShutdown([]io.Closer{kc, s})
    //notreached
}

func (as *appState) gracefulShutdown(cleanupFns []io.Closer) {

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    sig := <-sigChan

    fmt.Printf("terminating due to signal(%v)\n", sig)
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
    const MaxErrorCount = 10
    return func(err error) bool {

        errorCount++
        errorTime := time.Now()
        fmt.Printf("error %d/%d received(%v)\n", errorCount, MaxErrorCount, err)
        if !lastErrorTime.IsZero() && errorTime.Sub(lastErrorTime) > 60*time.Second {
            fmt.Printf("reset error counter\n")
            errorCount = 0
        }

        lastErrorTime = errorTime
        if errorCount < MaxErrorCount {
            return false
        }

        fmt.Printf("terminating process due to error count\n")
        _ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
        return true
    }
}
