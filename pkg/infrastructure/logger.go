package infrastructure

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})

	if os.Getenv("PE_DEBUG") != "true" {
		switch os.Getenv("LOG_LEVEL") {
		case "FATAL":
			log.SetLevel(logrus.FatalLevel)
		case "ERROR":
			log.SetLevel(logrus.ErrorLevel)
		case "WARN":
			log.SetLevel(logrus.WarnLevel)
		case "DEBUG":
			log.SetLevel(logrus.DebugLevel)
		case "INFO":
			log.SetLevel(logrus.InfoLevel)
		case "TRACE":
			log.SetLevel(logrus.TraceLevel)
		default:
			log.SetLevel(logrus.InfoLevel)
		}
	} else {
		log.SetLevel(logrus.TraceLevel)
	}

	return log
}
