package log

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is the global logger instance
var Logger zerolog.Logger

// InitLogger initializes the logger with caller information and sets the log level.
func InitLogger(logFileDir string, level zerolog.Level) error {
	// Ensure the log directory exists
	if err := os.MkdirAll(logFileDir, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("Failed to create log directory")
		return err
	}

	// Define the log file path
	logFilePath := filepath.Join(logFileDir, "overseer_api.log")

	// Lumberjack for log rotation
	logFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10, // megabytes
		MaxBackups: 3,  // number of backups to keep
		MaxAge:     28, // days to keep a log file
		Compress:   false,
	}

	// Create a multi-output writer for both console and file
	multi := io.MultiWriter(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    false,
			TimeFormat: "2006-01-02 15:04:05",
			FormatLevel: func(i interface{}) string {
				if ll, ok := i.(string); ok {
					switch strings.ToLower(ll) {
					case "trace":
						return "\033[34m" + strings.ToUpper(ll) + "\033[0m " // Blue
					case "debug":
						return "\033[35m" + strings.ToUpper(ll) + "\033[0m " // Magenta
					case "info":
						return "\033[36m" + strings.ToUpper(ll) + "\033[0m " // Cyan
					case "warn":
						return "\033[33m" + strings.ToUpper(ll) + "\033[0m " // Yellow
					case "error":
						return "\033[31m" + strings.ToUpper(ll) + "\033[0m " // Red
					case "fatal":
						return "\033[41m" + strings.ToUpper(ll) + "\033[0m " // Red background
					case "panic":
						return "\033[41m" + strings.ToUpper(ll) + "\033[0m " // Red background
					default:
						return strings.ToUpper(ll)
					}
				}
				return ""
			},
		},
		logFile,
	)

	// Initialize the logger with caller information and set the log level
	Logger = zerolog.New(multi).With().Timestamp().Caller().Logger().Level(level)
	return nil
}
