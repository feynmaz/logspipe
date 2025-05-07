package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	appName string = "logspipe"
	appEnv  string = "dev"

	port int = 8088
)

func ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func main() {
	zerolog.MessageFieldName = "_msg"
	stream := fmt.Sprintf("app=%s,env=%s", appName, appEnv)
	logger := zerolog.New(os.Stdout).With().
		Str("_stream", stream).
		Timestamp().
		Logger().
		Level(zerolog.Level(0))

	go produceLogs(logger)

	http.HandleFunc("/ready", ready)
	logger.Info().Msgf("Server starting on port %d..", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Fatal().Err(err).Msgf("server error")
	}
}

func produceLogs(logger zerolog.Logger) {
	var idx int = 100
	for {
		for i := range idx {
			logger.Info().Msgf("message %d", i)
			time.Sleep(200 * time.Millisecond)
		}
		time.Sleep(5 * time.Second)
	}
}
