package global

import (
	"fmt"
	"os"

	log "github.com/cihub/seelog"
)

func initLog() {
	logFile := os.Getenv("LOG_CONFIG_FILE")
	if logFile == "" {
		logFile = "../etc/" + Configs.Log
	}
	logger, err := log.LoggerFromConfigAsFile(logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load log config from %s failed: %s\n", Configs.Log, err)
		os.Exit(1)
	}

	log.ReplaceLogger(logger)
}
