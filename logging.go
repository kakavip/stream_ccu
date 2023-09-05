package main

import (
	"log"
	"os"
)

type Logging interface {
	info(msg string)
	error(msg string)
	warn(msg string)
}

type AppLogging struct {
	logger *log.Logger
}

func (logging *AppLogging) info(msg string) {
	logging.logger.Println("[INFO] " + msg)
}
func (logging *AppLogging) error(msg string) {
	logging.logger.Println("[ERROR] " + msg)
}
func (logging *AppLogging) warn(msg string) {
	logging.logger.Println("[ERROR] " + msg)
}

func NewLogger() Logging {
	return &AppLogging{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}
