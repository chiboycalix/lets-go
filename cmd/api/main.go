package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/chiboycalix/go-further/internal/data"
	"github.com/chiboycalix/go-further/internal/jsonlog"
	"github.com/chiboycalix/go-further/internal/mailer"
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

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}

	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	configuration configuration
	logger        *jsonlog.Logger
	models        data.Models
	mailer        mailer.Mailer
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

	flag.Float64Var(&config.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&config.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&config.limiter.enabled, "limiter-enabled", true, "Rate limiter enabled")

	flag.StringVar(&config.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&config.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&config.smtp.username, "smtp-username", "704954e8d8805b", "SMTP username")
	flag.StringVar(&config.smtp.password, "smtp-password", "293199d8d5c3a0", "SMTP password")
	flag.StringVar(&config.smtp.sender, "smtp-sender", "Greenlight <no-reply@igwechinonso77@gmail.com>", "SMTP sender")

	flag.Parse()
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
		mailer:        mailer.New(config.smtp.host, config.smtp.port, config.smtp.username, config.smtp.password, config.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
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
