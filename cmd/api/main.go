package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chiboycalix/go-further/internal/data"
	"github.com/chiboycalix/go-further/internal/jsonlog"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type configuration struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	configuration configuration
	// logger        *log.Logger
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var config configuration
	flag.IntVar(&config.port, "port", 1004, "API server port")
	flag.StringVar(&config.env, "env", "development", "Environment (development|staging|production)")
	// flag.StringVar(&config.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&config.db.dsn, "db-dsn", os.Getenv("MOVIE_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&config.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL Max open connection")
	flag.IntVar(&config.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQ Max Idle connection")
	flag.StringVar(&config.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQ Max Idle Time")
	flag.Parse()
	// logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(config)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	// psql -U postgres
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		configuration: config,
		logger:        logger,
		models:        data.NewModels(db),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//start http server
	logger.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  config.env,
	})
	err = server.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(config configuration) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.db.maxOpenConns)
	db.SetMaxIdleConns(config.db.maxIdleConns)

	duration, err := time.ParseDuration(config.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
