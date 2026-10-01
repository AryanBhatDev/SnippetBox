package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {

	port := flag.String("addr", ":4000", "server port")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	flag.Parse()

	app := application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", *port))

	err := http.ListenAndServe(*port, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}
