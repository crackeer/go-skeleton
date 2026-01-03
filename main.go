package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-skeleton/container"
	"go-skeleton/server"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "c", "./config/app.yaml", "config file")
	flag.Parse()

	if len(configPath) < 1 {
		fmt.Println("config file is required")
		os.Exit(-1)
		return
	}

	appConfig, err := container.Init(configPath)
	if err != nil {
		fmt.Printf("failed to init app: %s\n", err.Error())
		os.Exit(-1)
		return
	}

	errChan := make(chan error)

	go func() {
		err := server.Run(appConfig.Port)
		errChan <- err
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)

	select {
	case err := <-errChan:
		fmt.Printf("encounter error when starting server with [%s]\n", err.Error())
	case <-signalChan:
		fmt.Printf("received signal to shutdown, process will exit\n")
	}
}
