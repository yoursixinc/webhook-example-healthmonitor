package main

import (
	"fmt"
	"healthmonitor/internal/config"
	"healthmonitor/internal/db"
	"healthmonitor/internal/handler/api"
	"healthmonitor/internal/handler/webhook"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"

	_ "embed"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// If in terminal, set log level to debug and output to console
	if _, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS); err == nil {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseHost, cfg.DatabaseName)

	db, err := db.New(url)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	defer db.Close()

	webhookHandler, err := webhook.New(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create webhook handler")
	}

	apiHandler, err := api.New(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create API handler")
	}

	http.Handle("/", http.FileServer(http.Dir("./dist")))
	http.Handle("/api/", http.StripPrefix("/api", apiHandler))
	http.Handle("/webhook", webhookHandler)

	log.Info().Msg("Starting web server")
	err = http.ListenAndServe(cfg.ListenHTTP, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
