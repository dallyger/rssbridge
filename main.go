package main

import (
	"log/slog"
	"os"
	"slices"

	"vnbr.de/rssbridge/internal/router"
)

func main() {
	if slices.Contains(os.Args, "-vvv") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	router.Run()
}
