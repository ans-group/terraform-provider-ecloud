package logger

import "log"

type ProviderLogger struct{}

func (l *ProviderLogger) Error(msg string) {
	l.log("ERROR", msg)
}

func (l *ProviderLogger) Warn(msg string) {
	l.log("WARN", msg)
}

func (l *ProviderLogger) Info(msg string) {
	l.log("INFO", msg)
}

func (l *ProviderLogger) Debug(msg string) {
	l.log("DEBUG", msg)
}

func (l *ProviderLogger) Trace(msg string) {
	l.log("TRACE", msg)
}

func (l *ProviderLogger) log(level string, msg string) {
	log.Printf("[%s] %s", level, msg)
}
