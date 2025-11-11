package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Severity levels
const (
	SeverityDebug   = "DEBUG"
	SeverityInfo    = "INFO"
	SeverityWarning = "WARNING"
	SeverityError   = "ERROR"
	SeverityFatal   = "FATAL"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Severity  string                 `json:"severity"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// logStructured writes a structured log entry to stdout
func logStructured(severity, message string, metadata map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Severity:  severity,
		Message:   message,
		Service:   getServiceName(),
		Metadata:  metadata,
	}

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		// Fallback to standard logging if JSON marshaling fails
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}
	log.Println(string(jsonEntry))
}

// getServiceName returns the service name from environment or a default
func getServiceName() string {
	// K_SERVICE is used by Google Cloud Run
	if service := os.Getenv("K_SERVICE"); service != "" {
		return service
	}
	// Fallback to SERVICE_NAME or default
	if service := os.Getenv("SERVICE_NAME"); service != "" {
		return service
	}
	return "shroomp-service"
}

// Debug logs a debug message
func Debug(message string, metadata map[string]interface{}) {
	logStructured(SeverityDebug, message, metadata)
}

// Info logs an info message
func Info(message string, metadata map[string]interface{}) {
	logStructured(SeverityInfo, message, metadata)
}

// Warning logs a warning message
func Warning(message string, metadata map[string]interface{}) {
	logStructured(SeverityWarning, message, metadata)
}

// Error logs an error message
func Error(message string, metadata map[string]interface{}) {
	logStructured(SeverityError, message, metadata)
}

// Fatal logs a fatal message and exits
func Fatal(message string, metadata map[string]interface{}) {
	logStructured(SeverityFatal, message, metadata)
	os.Exit(1)
}
