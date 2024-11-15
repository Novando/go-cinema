package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"sync"
)

type Config struct {

	// ConsoleLoggingEnabled console logging
	ConsoleLoggingEnabled bool

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool

	// FileLoggingEnabled makes the framework log to a file.go
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool

	// Directory to log to to when filelogging is enabled
	Directory string

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int

	// MaxBackups the max number of rolled files to keep
	MaxBackups int

	// MaxAge the max age in days to keep a logfile
	MaxAge int

	// CallerSkip the number of directory hierarchy to be skipped
	CallerSkip int
}

type Logger struct {
	serviceLogger *zerolog.Logger
}

var (
	instance *Logger
	once     sync.Once
)

// InitZerolog, call it during cmd main.go
func InitZerolog(config Config) {
	once.Do(func() {
		instance = initializer(config)
	})
}

func initializer(config ...Config) *Logger {
	if len(config) == 0 {
		config = append(config, Config{
			ConsoleLoggingEnabled: true,
		})
	}
	var writers []io.Writer

	if config[0].ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config[0].FileLoggingEnabled {
		writers = append(writers, newRollingFile(config[0]))
	}
	if config[0].CallerSkip != 0 {
		zerolog.CallerSkipFrameCount = config[0].CallerSkip
	}
	mw := io.MultiWriter(writers...)

	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Bool("fileLogging", config[0].FileLoggingEnabled).
		Bool("jsonLogOutput", config[0].EncodeLogsAsJson).
		Str("logDirectory", config[0].Directory).
		Str("fileName", config[0].Filename).
		Int("maxSizeMB", config[0].MaxSize).
		Int("maxBackups", config[0].MaxBackups).
		Int("maxAgeInDays", config[0].MaxAge)

	return &Logger{
		serviceLogger: &logger,
	}
}

func newRollingFile(config Config) io.Writer {
	l := &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}

	return l
}

// Call the singleton function
func Call() *Logger {
	if instance == nil {
		fmt.Println("logger not initialized, Call InitZerolog() first.")
		return nil
	}
	return instance
}

// Errorf print formatted error message
func (l *Logger) Errorf(format string, a ...interface{}) {
	errs := fmt.Errorf(format, a...)
	l.serviceLogger.Error().Caller().Msgf(errs.Error())
}

// Infof Print log message with format
func (l *Logger) Infof(format string, a ...interface{}) {
	l.serviceLogger.Info().Msgf(format, a...)
}

// Info Print log message
func (l *Logger) Info(format string) {
	l.serviceLogger.Info().Msg(format)
}

// Fatalf
// Stop the app after invocation, and
// print Fatal message with format
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.serviceLogger.Fatal().Msgf(format, a...)
}

// Fatal
// Stop the app after invocation,
// and print the Fatal message
func (l *Logger) Fatal(format string) {
	l.serviceLogger.Fatal().Msg(format)
}

// Panic
// Stop the app after invocation,
// and print the Panic message
func (l *Logger) Panic(format string, a ...interface{}) {
	l.serviceLogger.Panic().Msgf(format, a...)
}
