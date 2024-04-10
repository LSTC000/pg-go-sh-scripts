package main

import (
	"log"
	"os"
	"os/signal"
	"pg-sh-scripts/internal/server"
	"syscall"
	_ "time/tzdata"
)

// @title           Bash Scripts
// @version         1.0.0
// @description     This is an API for running bash scripts

// @host      0.0.0.0:8000
// @BasePath  /api/v1

func main() {
	s := server.GetServer()

	go func() {
		if err := s.Run(); err != nil {
			log.Fatalf("Run server error: %v", err)
		}
	}()

	sqt := make(chan os.Signal, 1)
	signal.Notify(sqt, syscall.SIGTERM, syscall.SIGINT)
	<-sqt

	if err := s.Shutdown(); err != nil {
		log.Fatalf("Shutdown server error: %v", err)
	}
}
