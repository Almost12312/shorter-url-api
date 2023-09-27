package main

import (
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"short-url-api/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in load .env file!")
	}
}

func main() {
	// app conf
	conf := config.MustLoad()

	log := setupLogger(conf.Env)
	log.Info("Start url shortener!", slog.String("env", conf.Env))
	log.Debug("debug msg enabled!")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
