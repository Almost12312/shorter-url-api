package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"short-url-api/internal/http-server/handlers/redirect"
	deleteUrl "short-url-api/internal/http-server/handlers/url/delete"
	"short-url-api/internal/http-server/handlers/url/save"

	"short-url-api/internal/config"
	mwLogger "short-url-api/internal/http-server/middleware/logger"
	"short-url-api/internal/lib/logger/sl"
	"short-url-api/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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

	//logger
	log := setupLogger(conf.Env)
	log.Info("Start url shortener!", slog.String("env", conf.Env))
	log.Debug("debug msg enabled!")

	//storage
	storage, err := sqlite.New(conf.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	} else {
		log.Info("logger has been created!")
	}

	//id, err := storage.SaveUrl("https://github.com", "google")
	//if err != nil {
	//	log.Error("failed to save url", sl.Err(err))
	//}
	//
	//_, err = storage.SaveUrl("https://mail.com", "google")
	//if err != nil {
	//	log.Error("failed to save url", sl.Err(err))
	//}
	//
	//log.Info("saved url to git", slog.Int64("id", id))
	//
	//res, err := storage.GetUrl("google")
	//if err != nil {
	//	fmt.Printf("error is %w\n", err)
	//}
	//
	//_ = res
	//
	//ok, err := storage.DeleteUrl("google")
	//fmt.Println(ok)
	//if err != nil {
	//	fmt.Println(err)
	//}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	//router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/{alias}", redirect.New(log, storage))

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			conf.HttpServer.User: conf.HttpServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/delete/{alias}", deleteUrl.New(log, storage))
	})

	log.Info("starting server", slog.String("address", conf.Address))
	srv := &http.Server{
		Addr:         conf.Address,
		Handler:      router,
		ReadTimeout:  conf.HttpServer.Timeout,
		WriteTimeout: conf.HttpServer.Timeout,
		IdleTimeout:  conf.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed start server!")
	}

	log.Info("server stopped!")
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
