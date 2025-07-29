package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/leandroolgomes/golang-dependency-graph/component"
	"github.com/leandroolgomes/golang-dependency-graph/examples"
)

func main() {

	config := component.Define("config", new(examples.Config))
	appRoutes := component.Define("app_routes", new(examples.AppRoutes))
	httpServer := component.Define("http_server", new(examples.HttpServer), appRoutes.Key(), config.Key())


	components := map[string]*component.Component{
		config.Key():     config,
		appRoutes.Key():  appRoutes,
		httpServer.Key(): httpServer,
	}


	system := component.CreateSystem(components)

	fmt.Println("Starting system...")
	if err := system.Start(); err != nil {
		fmt.Printf("Failed to start system: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("System started successfully")


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)


	sig := <-sigChan
	fmt.Printf("%s signal received, shutting down...\n", sig)


	if err := system.Stop(); err != nil {
		fmt.Printf("Error during system shutdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("System stopped successfully")
}
