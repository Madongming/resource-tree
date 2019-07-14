package global

import (
	"fmt"
	"os"

	log "github.com/cihub/seelog"
)

func initLog() {
	logger, err := log.LoggerFromConfigAsFile("../etc/" + Configs.Log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load log config from %s failed: %s\n", Configs.Log, err)
		os.Exit(1)
	}

	log.ReplaceLogger(logger)
}
