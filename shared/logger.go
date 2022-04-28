package shared

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

// NewLogger returns a logger instance
func NewLogger(mode string) {
	log.SetLevel(log.InfoLevel)
	level, err := log.ParseLevel(mode)
	if err == nil {
		os.Exit(2)
	}

	log.SetLevel(level)

	// @TODO: Get log path and determine where to write out
	log.SetHandler(text.New(os.Stdout))
}
