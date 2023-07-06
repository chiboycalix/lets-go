package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type configuration struct {
	port int
	env  string
}

type application struct {
	configuration configuration
	logger        *log.Logger
}

func main() {
	var config configuration

	flag.IntVar(&config.port, "port", 1004, "API server port")
	flag.StringVar(&config.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		configuration: config,
		logger:        logger,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//start http server
	logger.Printf("Starting %s server on %s", config.env, server.Addr)
	err := server.ListenAndServe()
	logger.Fatal(err)
}

