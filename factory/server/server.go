package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vizucode/gokit/factory"
)

// server an instance for running services with factory.ApplicationFactory
type server struct {
	service factory.ServiceFactory
}

// Server is abstraction of application Server
type Server interface {
	// Run all server actives
	Run()
}

// New initiate server to running the application
func New(svc factory.ServiceFactory) Server {
	return &server{service: svc}
}

func (s *server) Run() {
	if len(s.service.GetApplications()) < 1 {
		log.Fatal(fmt.Errorf("no server/worker/broker running"))
	}

	err := make(chan error, len(s.service.GetApplications()))
	for _, app := range s.service.GetApplications() {
		go func(srv factory.ApplicationFactory) {
			defer func() {
				if r := recover(); r != nil {
					err <- fmt.Errorf("%s", r)
				}
			}()

			// run the server
			srv.Serve()
		}(app)
	}

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)
	signal.Notify(quitSignal, syscall.SIGTERM)

	log.Printf("Application %s ready to run\n", s.service.Name())

	select {
	case e := <-err:
		panic(e)
	case <-quitSignal:
		s.shutdown(quitSignal)
	}
}

func (s *server) shutdown(forceShutdown chan os.Signal) {
	log.Println("Gracefully shutdown... (press Ctrl+C or Cmd+C to force)")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)

		for _, srv := range s.service.GetApplications() {
			srv.Shutdown(ctx)
		}

		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-forceShutdown:
		log.Println("Force shutdown servers, workers and brokers")
		cancel()
	case <-ctx.Done():
		log.Println("Context Timeout")
	}
}
