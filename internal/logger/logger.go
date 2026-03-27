package logger

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Logger defines our custom structured logger
type Logger struct {
	Level string // e.g. "INFO", "WARN", "ERROR"
}

// New creates a new Logger instance
func New(level string) *Logger {
	return &Logger{Level: strings.ToUpper(level)}
}

// Info logs an informational message
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.log("INFO", msg, fields)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.log("WARN", msg, fields)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.log("ERROR", msg, fields)
}

// maskSensitiveData uses regex and string manipulation to hide passwords and tokens
func maskSensitiveData(input string) string {
	// Advanced String manipulation: mask anything that looks like a password or token in logs
	var builder strings.Builder
	builder.WriteString(input)
	
	result := builder.String()
	
	// Basic string replacement for common keys
	lowerResult := strings.ToLower(result)
	if strings.Contains(lowerResult, "password") || strings.Contains(lowerResult, "token") {
		// Use regex for more advanced replacement
		re := regexp.MustCompile(`(?i)(password|token)["\s:=]+([a-zA-Z0-9!@#$%^&*()_+]+)`)
		result = re.ReplaceAllString(result, `$1="***MASKED***"`)
	}
	
	return result
}

// log formats and prints the log message
func (l *Logger) log(level, msg string, fields map[string]interface{}) {
	now := time.Now().Format(time.RFC3339)
	
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[%s] %s - %s", now, level, msg))
	
	if len(fields) > 0 {
		builder.WriteString(" | ")
		for k, v := range fields {
			// Mask sensitive fields explicitly
			keyLower := strings.ToLower(k)
			valStr := fmt.Sprintf("%v", v)
			
			if strings.Contains(keyLower, "password") || strings.Contains(keyLower, "token") {
				valStr = "***MASKED***"
			} else {
				valStr = maskSensitiveData(valStr) // Mask any embedded secrets
			}
			
			builder.WriteString(fmt.Sprintf("%s=%s ", k, valStr))
		}
	}
	
	fmt.Println(builder.String())
}
