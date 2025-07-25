package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()
	
	// Configurar formato JSON para producción
	Log.SetFormatter(&logrus.JSONFormatter{})
	
	// Configurar nivel de log desde variable de entorno
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}
	
	Log.SetOutput(os.Stdout)
}

func GetLogger() *logrus.Logger {
	if Log == nil {
		Init()
	}
	return Log
}